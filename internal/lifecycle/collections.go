package lifecycle

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	"github.com/zephinzer/ebzbaybot/internal/collection"
	"github.com/zephinzer/ebzbaybot/internal/floorpricediff"
	"github.com/zephinzer/ebzbaybot/internal/storage"
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

type ScrapingOpts struct {
	Connection     *sql.DB
	ScrapeInterval time.Duration
	Storage        storage.Storage
}

func StartCollectionsScraping(opts ScrapingOpts) {
	everyInterval := time.NewTicker(opts.ScrapeInterval).C
	dataStorage := opts.Storage
	for {
		<-everyInterval
		currentTimestamp := time.Now().Unix()
		currentTimestampString := strconv.FormatInt(currentTimestamp, 10)

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
		databaseCollections := collection.Collections{}
		currentCollectionsMap := types.CollectionStorage{}
		for _, currentCollection := range currentCollections {
			collectionID := currentCollection.Collection
			currentCollectionsMap[collectionID] = types.Collection{
				Data:        currentCollection,
				LastUpdated: currentTimestampString,
			}
			localWhitelistCollection, err := collection.GetCollectionByIdentifier(collectionID)
			if err == nil { // only add to the database if the collection is whitelisted
				databaseCollections = append(
					databaseCollections,
					collection.New(currentCollection, collection.NewOpts{
						Aliases: localWhitelistCollection.Aliases,
						Label:   localWhitelistCollection.Label,
					}),
				)
			}
		}

		log.Infof("storing in databse...")
		if err := collection.Save(collection.SaveOpts{
			Collections: databaseCollections,
			Connection:  opts.Connection,
		}); err != nil {
			log.Warnf("failed to save collections to db: %s", err)
		}

		log.Infof("processing diffs...")
		newCollectionsMap := types.CollectionStorage{}
		newCollectionsCount := 0
		floorPriceChangesMap := types.CollectionDiffStorage{}
		floorPriceChangesCount := 0
		databaseFloorPriceDiffs := floorpricediff.FloorPriceDiffs{}
		for currentCollectionKey, currentCollection := range currentCollectionsMap {
			_, exists := previousCollectionsMap[currentCollectionKey]
			if !exists {
				newCollectionsCount += 1
				newCollectionsMap[currentCollectionKey] = currentCollection
				continue
			}
			previousPrice := previousCollectionsMap[currentCollectionKey].Data.FloorPrice
			currentPrice := currentCollection.Data.FloorPrice
			if previousPrice != currentPrice {
				floorPriceChangesCount += 1
				floorPriceChangesMap[currentCollectionKey] = types.CollectionDiff{
					Current:     currentCollection,
					Previous:    previousCollectionsMap[currentCollectionKey],
					LastUpdated: currentTimestampString,
				}
				databaseFloorPriceDiffs = append(
					databaseFloorPriceDiffs,
					floorpricediff.FloorPriceDiff{
						CollectionID:  currentCollectionKey,
						PreviousPrice: previousPrice,
						CurrentPrice:  currentPrice,
						LastUpdated:   time.Unix(currentTimestamp, 0),
					},
				)
			}
		}

		log.Infof("storing %v new collections changes...", newCollectionsCount)
		newCollectionsJSON, err := json.Marshal(newCollectionsMap)
		if err != nil {
			log.Warnf("failed to marshal collection changes into json for storage: %s", err)
			continue
		}
		dataStorage.Set(StorageKeyNewCollections, newCollectionsJSON)
		dataStorage.Set(StorageKeyNewCollections+"_last_updated", []byte(currentTimestampString))

		log.Infof("storing %v floor price changes...", floorPriceChangesCount)
		floorPriceChangesJSON, err := json.Marshal(floorPriceChangesMap)
		if err != nil {
			log.Warnf("failed to marshal collection price changes into json for storage: %s", err)
			continue
		}
		dataStorage.Set(StorageKeyFloorPriceChanges, floorPriceChangesJSON)
		dataStorage.Set(StorageKeyFloorPriceChanges+"_last_updated", []byte(currentTimestampString))
		log.Infof("storing floor price diffs into db...")
		if err := floorpricediff.Save(floorpricediff.SaveOpts{
			Connection:      opts.Connection,
			FloorPriceDiffs: databaseFloorPriceDiffs,
		}); err != nil {
			log.Warnf("failed to save floor price diffs to db: %s", err)
		}

		log.Infof("storing current collection...")
		currentCollectionsMapJSON, err := json.Marshal(currentCollectionsMap)
		if err != nil {
			log.Warnf("failed to marshal collections into json for storage: %s", err)
			continue
		}
		dataStorage.Set(StorageKeyCollections, currentCollectionsMapJSON)
		dataStorage.Set(StorageKeyCollections+"_last_updated", []byte(currentTimestampString))
	}
}
