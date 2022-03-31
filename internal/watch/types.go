package watch

import "time"

type ChannelWatches []ChannelWatch

type ChannelWatch struct {
	ChatID       string    `json:"chatID"`
	CollectionID string    `json:"collectionID"`
	LastUpdated  time.Time `json:"lastUpdated"`
}

type Watches []Watch

type Watch struct {
	ChatID       int64     `json:"chatID"`
	CollectionID string    `json:"collectionID"`
	LastUpdated  time.Time `json:"lastUpdated"`
}
