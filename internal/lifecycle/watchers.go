package lifecycle

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/collection"
	"github.com/zephinzer/ebzbaybot/internal/floorpricediff"
	"github.com/zephinzer/ebzbaybot/internal/storage"
	"github.com/zephinzer/ebzbaybot/internal/types"
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
		log.Infof("loading watches from mem...")
		watchesJSON, _ := opts.Storage.Get("watches")

		log.Infof("loading watches from db...")
		databaseWatches, err := watch.Load(watch.LoadOpts{
			Connection: opts.Connection,
		})
		if err != nil {
			log.Warnf("failed to load watches from db: %s", err)
		}
		log.Infof("%v watches loaded from db", len(databaseWatches))

		// get floor price changes
		log.Infof("loading floor price diffs from mem...")
		floorPriceChanges, _ := opts.Storage.Get(StorageKeyFloorPriceChanges)
		floorPriceChangesMap := types.CollectionDiffStorage{}
		json.Unmarshal(floorPriceChanges, &floorPriceChangesMap)

		log.Infof("loading floor price diffs from db...")
		databaseFloorPriceDiffs, err := floorpricediff.Load(floorpricediff.LoadOpts{
			Connection: opts.Connection,
		})
		if err != nil {
			log.Warnf("failed to load watches from db: %s", err)
		}
		log.Infof("%v floor price differences loaded from db", len(databaseFloorPriceDiffs))

		watchesMap := types.WatchStorage{}
		json.Unmarshal(watchesJSON, &watchesMap)

		for userID, watch := range watchesMap {
			collectionsCount := 0
			for collectionID, lastUpdated := range watch.CollectionMap {
				collectionsCount += 1
				watchLastUpdated, _ := strconv.ParseInt(lastUpdated, 10, 64)
				floorPriceLastUpdated, _ := strconv.ParseInt(floorPriceChangesMap[collectionID].LastUpdated, 10, 64)
				isUpdateDue := floorPriceLastUpdated > watchLastUpdated

				if isUpdateDue {
					previousFloorPrice := floorPriceChangesMap[collectionID].Previous.Data.FloorPrice
					previousFloorPriceFloat, _ := strconv.ParseFloat(previousFloorPrice, 64)
					currentFloorPrice := floorPriceChangesMap[collectionID].Current.Data.FloorPrice
					currentFloorPriceFloat, _ := strconv.ParseFloat(currentFloorPrice, 64)

					// send
					chatID, _ := strconv.ParseInt(userID, 10, 64)
					collectionInstance, _ := collection.GetCollectionByIdentifier(collectionID)
					directionText := "🔻"
					if currentFloorPriceFloat > previousFloorPriceFloat {
						directionText = "🔼"
					}
					msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
						"🚨%s Floor price of [%s collection](https://app.ebisusbay.com/collection/%s) has changed from _%s_ CRO to *%s* CRO",
						directionText,
						collectionInstance.Label,
						collectionInstance.ID,
						previousFloorPrice,
						currentFloorPrice,
					))
					msg.ParseMode = "markdown"
					opts.Bot.Send(msg)

					currentTimestamp := strconv.FormatInt(time.Now().Unix(), 10)
					watchesMap[userID].CollectionMap[collectionID] = currentTimestamp
				}
			}
			log.Infof("user[%s] is watching %v collections", userID, collectionsCount)
		}
		watchesJSON, _ = json.Marshal(watchesMap)
		opts.Storage.Set("watches", watchesJSON)
	}

}
