package lifecycle

import (
	"database/sql"
	"time"

	"github.com/zephinzer/ebzbaybot/internal/collection"
	"github.com/zephinzer/ebzbaybot/internal/floorpricediff"
	"github.com/zephinzer/ebzbaybot/internal/storage"
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
	for {
		<-everyInterval
		currentTimestamp := time.Now().Unix()

		// retrieve existing values
		previousCollections, err := collection.LoadAll(collection.LoadOpts{
			Connection: opts.Connection,
		})
		if err != nil {
			log.Warnf("failed to retrieve collections: %s", err)
		}
		previousCollectionsMap := map[string]collection.Collection{}
		for _, previousDatabaseCollection := range previousCollections {
			previousCollectionsMap[previousDatabaseCollection.ID] = previousDatabaseCollection
		}

		// retrieve current values
		currentEbzbayCollections, err := ebzbay.GetCollections()
		if err != nil {
			log.Warnf("failed to get collections data: %s", err)
			continue
		}
		currentCollections := collection.Collections{}
		for _, currentCollection := range currentEbzbayCollections {
			collectionID := currentCollection.Collection
			localWhitelistCollection, err := collection.GetCollectionByIdentifier(collectionID)
			if err == nil { // only add to the database if the collection is whitelisted
				currentCollections = append(
					currentCollections,
					collection.New(currentCollection, collection.NewOpts{
						Aliases: localWhitelistCollection.Aliases,
						Label:   localWhitelistCollection.Label,
					}),
				)
			}
		}
		log.Infof("processing %v/%v collections from api", len(currentCollections), len(currentEbzbayCollections))
		currentCollectionsMap := map[string]collection.Collection{}
		for _, currentDatabaseCollection := range currentCollections {
			currentCollectionsMap[currentDatabaseCollection.ID] = currentDatabaseCollection
		}

		// process differences
		newDatabaseCollectionsCount := 0
		databaseFloorPriceDiffs := floorpricediff.FloorPriceDiffs{}
		for currentCollectionKey, currentCollection := range currentCollectionsMap {
			if _, exists := previousCollectionsMap[currentCollectionKey]; !exists {
				newDatabaseCollectionsCount += 1
				continue
			}
			previousPrice := previousCollectionsMap[currentCollectionKey].FloorPrice
			currentPrice := currentCollection.FloorPrice
			if previousPrice != currentPrice {
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
		log.Infof("evaluated %v changes in floor prices", len(databaseFloorPriceDiffs))

		// save floor prices
		if err := floorpricediff.Save(floorpricediff.SaveOpts{
			Connection:      opts.Connection,
			FloorPriceDiffs: databaseFloorPriceDiffs,
		}); err != nil {
			log.Warnf("failed to save floor price diffs to db: %s", err)
		}

		// save current collections
		log.Infof("saving %v collections to db...", len(currentCollections))
		if err := collection.Save(collection.SaveOpts{
			Collections: currentCollections,
			Connection:  opts.Connection,
		}); err != nil {
			log.Warnf("failed to save collections to db: %s", err)
		}
	}
}
