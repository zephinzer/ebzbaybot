package floorpricediff

import "time"

type FloorPriceDiffs []FloorPriceDiff

type FloorPriceDiff struct {
	CollectionID  string    `json:"collectionID"`
	PreviousPrice string    `json:"previousPrice"`
	CurrentPrice  string    `json:"currentPrice"`
	LastUpdated   time.Time `json:"lastUpdated"`
}
