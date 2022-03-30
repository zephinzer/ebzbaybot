package lifecycle

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/collection"
	"github.com/zephinzer/ebzbaybot/internal/constants"
	"github.com/zephinzer/ebzbaybot/internal/floorpricediff"
	"github.com/zephinzer/ebzbaybot/internal/storage"
	"github.com/zephinzer/ebzbaybot/internal/utils/log"
	"github.com/zephinzer/ebzbaybot/internal/watch"
)

type WatchingOpts struct {
	Bot        *tgbotapi.BotAPI
	Connection *sql.DB
	Storage    storage.Storage
}

func StartUpdatingWatchers(opts WatchingOpts) error {
	everyInterval := time.NewTicker(5 * time.Second).C
	for {
		<-everyInterval
		// load watches
		databaseWatches, err := watch.Load(watch.LoadOpts{
			Connection: opts.Connection,
		})
		if err != nil {
			log.Warnf("failed to load watches: %s", err)
		}
		log.Infof("loaded %v watches", len(databaseWatches))

		// load floor price changes
		databaseFloorPriceDiffs, err := floorpricediff.Load(floorpricediff.LoadOpts{
			Connection: opts.Connection,
		})
		if err != nil {
			log.Warnf("failed to load watches from db: %s", err)
		}
		databaseFloorPriceDiffsMap := map[string]floorpricediff.FloorPriceDiff{}
		for _, databaseFloorPriceDiff := range databaseFloorPriceDiffs {
			databaseFloorPriceDiffsMap[databaseFloorPriceDiff.CollectionID] = databaseFloorPriceDiff
		}
		log.Infof("loaded %v floor price differences", len(databaseFloorPriceDiffs))

		// go through watches and check if last updated floor price is earlier than last user updatred
		databaseUpdatedWatches := watch.Watches{}
		for _, databaseWatch := range databaseWatches {
			collectionID := databaseWatch.CollectionID
			userLastUpdatedAt := databaseWatch.LastUpdated
			floorPriceLastUpdatedAt := databaseFloorPriceDiffsMap[collectionID].LastUpdated
			if floorPriceLastUpdatedAt.After(userLastUpdatedAt) {
				collectionInstance, _ := collection.GetCollectionByIdentifier(collectionID)

				// trigger user update
				directionSymbol := constants.UserTextPriceDown
				directionText := "down"
				previousFloorPrice := databaseFloorPriceDiffsMap[collectionID].PreviousPrice
				previousFloorPriceFloat, _ := strconv.ParseFloat(previousFloorPrice, 64)
				currentFloorPrice := databaseFloorPriceDiffsMap[collectionID].CurrentPrice
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
				opts.Bot.Send(msg)

				databaseWatch.LastUpdated = time.Now()
				databaseUpdatedWatches = append(databaseUpdatedWatches, databaseWatch)
			}
		}
		if err := watch.Save(watch.SaveOpts{
			Connection: opts.Connection,
			Watches:    databaseUpdatedWatches,
		}); err != nil {
			log.Warnf("failed to save %v watches to database: %s", len(databaseUpdatedWatches), err)
		} else {
			log.Infof("updated %v watches", len(databaseUpdatedWatches))
		}
	}
}
