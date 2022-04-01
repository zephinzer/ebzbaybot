package floorpricediff

import "time"

type FloorPriceDiffs []FloorPriceDiff

type FloorPriceDiff struct {
	CollectionID  string    `json:"collectionID"`
	PreviousPrice string    `json:"previousPrice"`
	CurrentPrice  string    `json:"currentPrice"`
	ListingID     *string   `json:"listingID"`
	ImageURL      *string   `json:"imageURL"`
	Edition       *string   `json:"edition"`
	Score         *string   `json:"score"`
	LastUpdated   time.Time `json:"lastUpdated"`
}
