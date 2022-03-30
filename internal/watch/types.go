package watch

import "time"

type Watches []Watch

type Watch struct {
	ChatID       string    `json:"chatID"`
	CollectionID string    `json:"collectionID"`
	LastUpdated  time.Time `json:"lastUpdated"`
}
