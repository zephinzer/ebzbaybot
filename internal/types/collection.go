package types

import "github.com/zephinzer/ebzbaybot/pkg/ebzbay"

type CollectionDiffStorage map[string]CollectionDiff

type CollectionDiff struct {
	Current     Collection `json:"current"`
	Previous    Collection `json:"previous"`
	LastUpdated string     `json:"lastUpdated"`
}

type CollectionStorage map[string]Collection

type Collection struct {
	Data        ebzbay.Collection `json:"data"`
	LastUpdated string            `json:"lastUpdated"`
}
