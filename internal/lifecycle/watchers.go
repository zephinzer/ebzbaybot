package lifecycle

import (
	"database/sql"
	"fmt"
	"path"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/collection"
	"github.com/zephinzer/ebzbaybot/internal/constants"
	"github.com/zephinzer/ebzbaybot/internal/floorpricediff"
	"github.com/zephinzer/ebzbaybot/internal/telegram/handlers"
	"github.com/zephinzer/ebzbaybot/internal/utils/log"
	"github.com/zephinzer/ebzbaybot/internal/watch"
)

type WatchingOpts struct {
	Bot        *tgbotapi.BotAPI
	Connection *sql.DB
}

func StartUpdatingWatchers(opts WatchingOpts) error {
	everyInterval := time.NewTicker(5 * time.Second).C
	for {
		<-everyInterval
		// load watches
		chatWatches, err := watch.Load(watch.LoadOpts{
			Connection: opts.Connection,
		})
		if err != nil {
			log.Warnf("failed to load chat watches: %s", err)
		}
		log.Infof("loaded %v chat watches", len(chatWatches))

		// load floor price changes
		floorPriceDiffs, err := floorpricediff.Load(floorpricediff.LoadOpts{
			Connection: opts.Connection,
		})
		if err != nil {
			log.Warnf("failed to load chat watches from db: %s", err)
		}
		floorPriceDiffsMap := map[string]floorpricediff.FloorPriceDiff{}
		for _, databaseFloorPriceDiff := range floorPriceDiffs {
			floorPriceDiffsMap[databaseFloorPriceDiff.CollectionID] = databaseFloorPriceDiff
		}
		log.Infof("loaded %v floor price differences", len(floorPriceDiffs))

		// go through watches and check if last updated floor price is earlier than last user updatred
		updatedChatWatches := watch.Watches{}
		for _, databaseWatch := range chatWatches {
			collectionID := databaseWatch.CollectionID
			userLastUpdatedAt := databaseWatch.LastUpdated
			floorPriceLastUpdatedAt := floorPriceDiffsMap[collectionID].LastUpdated
			if floorPriceLastUpdatedAt.After(userLastUpdatedAt) {
				collectionInstance, _ := collection.GetCollectionByIdentifier(collectionID)

				// trigger user update
				directionSymbol := constants.UserTextPriceDown
				directionText := "down"
				previousFloorPrice := floorPriceDiffsMap[collectionID].PreviousPrice
				previousFloorPriceFloat, _ := strconv.ParseFloat(previousFloorPrice, 64)
				currentFloorPrice := floorPriceDiffsMap[collectionID].CurrentPrice
				currentFloorPriceFloat, _ := strconv.ParseFloat(currentFloorPrice, 64)
				if currentFloorPriceFloat > previousFloorPriceFloat {
					directionSymbol = constants.UserTextPriceUp
					directionText = "up"
				}
				log.Infof("triggering floor price change message to chat[%v]...", databaseWatch.ChatID)
				msg := tgbotapi.NewMessage(databaseWatch.ChatID, fmt.Sprintf(
					"🚨%s [%s](https://app.ebisusbay.com/collection/%s) FP: *%s* CRO (%s from _%s_ CRO)",
					directionSymbol,
					collectionInstance.Label,
					collectionInstance.ID,
					currentFloorPrice,
					directionText,
					previousFloorPrice,
				))
				msg.ParseMode = "markdown"
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData(
							"SHOW COLLECTION INFO",
							path.Join(
								handlers.CALLBACK_LIST_GET_NO_DELETE,
								collectionInstance.ID,
							),
						),
					),
				)
				opts.Bot.Send(msg)

				databaseWatch.LastUpdated = time.Now()
				updatedChatWatches = append(updatedChatWatches, databaseWatch)
			}
		}
		if err := watch.Save(watch.SaveOpts{
			Connection: opts.Connection,
			Watches:    updatedChatWatches,
		}); err != nil {
			log.Warnf("failed to save %v chat watches to database: %s", len(updatedChatWatches), err)
		} else {
			log.Infof("updated %v chat watches", len(updatedChatWatches))
		}
	}
}
