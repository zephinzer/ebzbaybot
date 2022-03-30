package collection

import "time"

type Collections []Collection

type Collection struct {
	ID               string    `json:"id"`
	Label            string    `json:"label"`
	Aliases          string    `json:"aliases"`
	AverageSalePrice string    `json:"averageSalePrice"`
	FloorPrice       string    `json:"floorPrice"`
	NumberActive     string    `json:"numberActive"`
	NumberOfSales    string    `json:"numberOfSales"`
	TotalRoyalties   string    `json:"totalRoyalties"`
	TotalVolume      string    `json:"totalVolume"`
	LastUpdated      time.Time `json:"lastUpdated"`
}
