package lifecycle

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/collection"
	"github.com/zephinzer/ebzbaybot/internal/constants"
	"github.com/zephinzer/ebzbaybot/internal/floorpricediff"
	"github.com/zephinzer/ebzbaybot/internal/utils/log"
	"github.com/zephinzer/ebzbaybot/internal/watch"
)

func StartUpdatingChannelWatchers(opts WatchingOpts) error {
	everyInterval := time.NewTicker(3 * time.Second).C
	for {
		<-everyInterval
		// load watches
		channelWatches, err := watch.LoadChannel(watch.LoadChannelOpts{
			Connection: opts.Connection,
		})
		if err != nil {
			log.Warnf("failed to channel load watches: %s", err)
		}
		log.Infof("loaded %v channel watches", len(channelWatches))

		// load floor price changes
		floorPriceDiffs, err := floorpricediff.Load(floorpricediff.LoadOpts{
			Connection: opts.Connection,
		})
		if err != nil {
			log.Warnf("failed to load channel watches from db: %s", err)
		}
		floorPriceDiffsMap := map[string]floorpricediff.FloorPriceDiff{}
		for _, databaseFloorPriceDiff := range floorPriceDiffs {
			floorPriceDiffsMap[databaseFloorPriceDiff.CollectionID] = databaseFloorPriceDiff
		}
		log.Infof("loaded %v floor price differences", len(floorPriceDiffs))

		// go through watches and check if last updated floor price is earlier than last user updatred
		updatedChannelWatches := watch.ChannelWatches{}
		for _, channelWatch := range channelWatches {
			collectionID := channelWatch.CollectionID
			floorPriceDiff := floorPriceDiffsMap[collectionID]
			userLastUpdatedAt := channelWatch.LastUpdated
			floorPriceLastUpdatedAt := floorPriceDiff.LastUpdated
			if floorPriceLastUpdatedAt.After(userLastUpdatedAt) {
				chatID := channelWatch.ChatID
				// this is because of a previous bug, if a wrong chat id was saved,
				// just skip it
				if len(chatID) == 0 {
					continue
				}
				if chatID[0] != '-' {
					chatID = "@" + chatID
				}

				collectionInstance, _ := collection.GetCollectionByIdentifier(collectionID)

				// trigger user update
				directionSymbol := constants.UserTextPriceDown
				directionText := "down"
				previousFloorPrice := floorPriceDiff.PreviousPrice
				currentFloorPrice := floorPriceDiff.CurrentPrice

				// this was added because of a very weird and flaky bug
				// that i cannot catch where a floor price diff was added
				// even though the prices is the same
				if strings.Compare(previousFloorPrice, currentFloorPrice) == 0 {
					continue
				}

				previousFloorPriceFloat, _ := strconv.ParseFloat(previousFloorPrice, 64)
				currentFloorPriceFloat, _ := strconv.ParseFloat(currentFloorPrice, 64)
				if currentFloorPriceFloat > previousFloorPriceFloat {
					directionSymbol = constants.UserTextPriceUp
					directionText = "up"
				}

				log.Infof("triggering floor price change message to channel[%s]...", chatID)

				additionalKeys := []string{"*Edition*"}
				additionalValues := []string{"*" + *floorPriceDiff.Edition + "*"}
				if floorPriceDiff.Rank != nil && *floorPriceDiff.Rank != "0" {
					additionalKeys = append(additionalKeys, "`Rank`")
					additionalValues = append(additionalValues, "`"+*floorPriceDiff.Rank+"`")
				}
				if floorPriceDiff.Score != nil && *floorPriceDiff.Score != "0.00" {
					additionalKeys = append(additionalKeys, "_Score_")
					additionalValues = append(additionalValues, "_"+*floorPriceDiff.Score+"_")
				}
				additionalProperties := fmt.Sprintf(
					"%s\n%s",
					strings.Join(additionalKeys, AdditionalPropertiesDelimiter),
					strings.Join(additionalValues, AdditionalPropertiesDelimiter),
				)

				msg := tgbotapi.NewMessageToChannel(chatID, fmt.Sprintf(
					"[🚨](%s)%s [%s](https://app.ebisusbay.com/collection/%s) FP: *%s* CRO (%s from _%s_ CRO)\n\n%s",
					*floorPriceDiff.ImageURL,
					directionSymbol,
					collectionInstance.Label,
					collectionInstance.ID,
					currentFloorPrice,
					directionText,
					previousFloorPrice,
					additionalProperties,
				))
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonURL(
							"🔗 BROWSER",
							fmt.Sprintf("https://app.ebisusbay.com/listing/%s", *floorPriceDiff.ListingID),
						),
						tgbotapi.NewInlineKeyboardButtonURL(
							"📲 METAMASK",
							fmt.Sprintf("https://metamask.app.link/dapp/app.ebisusbay.com/listing/%s", *floorPriceDiff.ListingID),
						),
					),
				)
				msg.ParseMode = "markdown"
				if _, err := opts.Bot.Send(msg); err != nil {
					log.Warnf("failed to send message to chat[%s]: %s", channelWatch.ChatID, err)
				}

				channelWatch.LastUpdated = time.Now()
				updatedChannelWatches = append(updatedChannelWatches, channelWatch)
			}
		}
		if err := watch.SaveChannel(watch.SaveChannelOpts{
			Connection: opts.Connection,
			Watches:    updatedChannelWatches,
		}); err != nil {
			log.Warnf("failed to save %v channel watches to database: %s", len(updatedChannelWatches), err)
		} else {
			log.Infof("updated %v channel watches", len(updatedChannelWatches))
		}
	}
}
