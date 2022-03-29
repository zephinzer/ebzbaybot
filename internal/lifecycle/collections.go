package lifecycle

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/zephinzer/ebzbaybot/internal/types"
	"github.com/zephinzer/ebzbaybot/internal/utils/log"
	"github.com/zephinzer/ebzbaybot/pkg/ebzbay"
)

const (
	StorageKeyCollectionUpdatedTimestamp = "collection_updated_"
	StorageKeyCollectionsLastUpdated     = "collections_last_updated"
	StorageKeyFloorPriceChanges          = "collections_floor_price_changes"
	StorageKeyNewCollections             = "collections_changes"
	StorageKeyCollections                = "collections"
)

func StartCollectionsScraping(opts ScrapingOpts) {
	everyInterval := time.NewTicker(opts.ScrapeInterval).C
	dataStorage := opts.Storage
	for {
		<-everyInterval
		currentTimestamp := strconv.FormatInt(time.Now().Unix(), 10)

		log.Infof("retrieving stored collections...")
		previousCollectionsMapJSON, _ := dataStorage.Get(StorageKeyCollections)
		previousCollectionsMap := types.CollectionStorage{}
		json.Unmarshal(previousCollectionsMapJSON, &previousCollectionsMap)

		log.Infof("retrieving updated collections...")
		currentCollections, err := ebzbay.GetCollections()
		if err != nil {
			log.Warnf("failed to get collections data: %s", err)
			continue
		}
		currentCollectionsMap := types.CollectionStorage{}
		for _, currentCollection := range currentCollections {
			currentCollectionsMap[currentCollection.Collection] = types.Collection{
				Data:        currentCollection,
				LastUpdated: currentTimestamp,
			}
		}

		log.Infof("processing diffs...")
		newCollectionsMap := types.CollectionStorage{}
		newCollectionsCount := 0
		floorPriceChangesMap := types.CollectionDiffStorage{}
		floorPriceChangesCount := 0
		for currentCollectionKey, currentCollection := range currentCollectionsMap {
			_, exists := previousCollectionsMap[currentCollectionKey]
			if !exists {
				newCollectionsCount += 1
				newCollectionsMap[currentCollectionKey] = currentCollection
				continue
			}
			if previousCollectionsMap[currentCollectionKey].Data.FloorPrice != currentCollection.Data.FloorPrice {
				floorPriceChangesCount += 1
				floorPriceChangesMap[currentCollectionKey] = types.CollectionDiff{
					Current:     currentCollection,
					Previous:    previousCollectionsMap[currentCollectionKey],
					LastUpdated: currentTimestamp,
				}
			}
		}

		log.Infof("storing %v new collections changes...", newCollectionsCount)
		newCollectionsJSON, err := json.Marshal(newCollectionsMap)
		if err != nil {
			log.Warnf("failed to marshal collection changes into json for storage: %s", err)
			continue
		}
		dataStorage.Set(StorageKeyNewCollections, newCollectionsJSON)
		dataStorage.Set(StorageKeyNewCollections+"_last_updated", []byte(currentTimestamp))

		log.Infof("storing %v floor price changes...", floorPriceChangesCount)
		floorPriceChangesJSON, err := json.Marshal(floorPriceChangesMap)
		if err != nil {
			log.Warnf("failed to marshal collection price changes into json for storage: %s", err)
			continue
		}
		dataStorage.Set(StorageKeyFloorPriceChanges, floorPriceChangesJSON)
		dataStorage.Set(StorageKeyFloorPriceChanges+"_last_updated", []byte(currentTimestamp))

		log.Infof("storing current collection...")
		currentCollectionsMapJSON, err := json.Marshal(currentCollectionsMap)
		if err != nil {
			log.Warnf("failed to marshal collections into json for storage: %s", err)
			continue
		}
		dataStorage.Set(StorageKeyCollections, currentCollectionsMapJSON)
		dataStorage.Set(StorageKeyCollections+"_last_updated", []byte(currentTimestamp))
	}
}
