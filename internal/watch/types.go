package watch

import "time"

type Watches []Watch

type Watch struct {
	ChatID       int64     `json:"chatID"`
	CollectionID string    `json:"collectionID"`
	LastUpdated  time.Time `json:"lastUpdated"`
}
