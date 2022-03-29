package lifecycle

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/collection"
	"github.com/zephinzer/ebzbaybot/internal/types"
	"github.com/zephinzer/ebzbaybot/internal/utils/log"
)

type Watcher struct {
	LastUpdated string `json:"lastUpdated"`
}

func StartUpdatingWatchers(opts WatchingOpts) error {
	everyInterval := time.NewTicker(5 * time.Second).C
	for {
		<-everyInterval
		watchesJSON, _ := opts.Storage.Get("watches")

		// get floor price changes
		floorPriceChanges, _ := opts.Storage.Get(StorageKeyFloorPriceChanges)
		floorPriceChangesMap := types.CollectionDiffStorage{}
		json.Unmarshal(floorPriceChanges, &floorPriceChangesMap)

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
						collectionInstance.Name,
						collectionInstance.Address,
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
