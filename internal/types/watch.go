package types

type WatchStorage map[string]Watch

type Watch struct {
	// CollectionMap is a key-value pair of collection ID and last updated unix time
	CollectionMap map[string]string `json:"collectionMap"`
}
