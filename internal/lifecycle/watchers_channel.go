package lifecycle

import (
	"fmt"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/collection"
	"github.com/zephinzer/ebzbaybot/internal/constants"
	"github.com/zephinzer/ebzbaybot/internal/floorpricediff"
	"github.com/zephinzer/ebzbaybot/internal/utils/log"
	"github.com/zephinzer/ebzbaybot/internal/watch"
)

func StartUpdatingChannelWatchers(opts WatchingOpts) error {
	everyInterval := time.NewTicker(5 * time.Second).C
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
			userLastUpdatedAt := channelWatch.LastUpdated
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
				log.Infof("triggering floor price change message to channel[%s]...", channelWatch.ChatID)
				msg := tgbotapi.NewMessageToChannel("@"+channelWatch.ChatID, fmt.Sprintf(
					"🚨%s [%s](https://app.ebisusbay.com/collection/%s) FP: *%s* CRO (%s from _%s_ CRO)",
					directionSymbol,
					collectionInstance.Label,
					collectionInstance.ID,
					currentFloorPrice,
					directionText,
					previousFloorPrice,
				))
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
