package watcher

import "time"

// Watchlist should be a map of chat IDs to collections watched
type Watchlist map[string][]string

type Config struct {
	LastUpdated time.Time
}
