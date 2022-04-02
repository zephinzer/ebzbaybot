package lifecycle

import (
	"database/sql"
	"time"

	"github.com/zephinzer/ebzbaybot/internal/collection"
	"github.com/zephinzer/ebzbaybot/internal/floorpricediff"
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
		currentCollectionsIndexMap := map[string]int{}
		for index, currentDatabaseCollection := range currentCollections {
			currentCollectionsMap[currentDatabaseCollection.ID] = currentDatabaseCollection
			currentCollectionsIndexMap[currentDatabaseCollection.ID] = index
		}

		// process differences
		newDatabaseCollectionsCount := 0
		floorPriceDiffs := floorpricediff.FloorPriceDiffs{}
		for currentCollectionKey, currentCollection := range currentCollectionsMap {
			if _, exists := previousCollectionsMap[currentCollectionKey]; !exists {
				newDatabaseCollectionsCount += 1
				continue
			}
			previousPrice := previousCollectionsMap[currentCollectionKey].FloorPrice
			currentPrice := currentCollection.FloorPrice
			if previousPrice != currentPrice {
				floorPriceDiffs = append(
					floorPriceDiffs,
					floorpricediff.FloorPriceDiff{
						CollectionID:  currentCollectionKey,
						PreviousPrice: previousPrice,
						CurrentPrice:  currentPrice,
						LastUpdated:   time.Unix(currentTimestamp, 0),
					},
				)
			}
		}
		log.Infof("evaluated %v changes in floor prices", len(floorPriceDiffs))

		for i := 0; i < len(floorPriceDiffs); i++ {
			floorCollectionStats := ebzbay.GetCollectionStats(floorPriceDiffs[i].CollectionID)
			floorPriceDiffs[i].ListingID = &floorCollectionStats.FloorListingID
			floorPriceDiffs[i].ImageURL = &floorCollectionStats.FloorImageURL
			floorPriceDiffs[i].Edition = &floorCollectionStats.FloorEdition
			floorPriceDiffs[i].Score = &floorCollectionStats.FloorScore
			floorPriceDiffs[i].Rank = &floorCollectionStats.FloorRank
		}

		// save floor prices
		if err := floorpricediff.Save(floorpricediff.SaveOpts{
			Connection:      opts.Connection,
			FloorPriceDiffs: floorPriceDiffs,
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
