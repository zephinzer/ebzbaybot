package stats

import (
	"database/sql"
	"fmt"
	"runtime"
	"time"

	"github.com/zephinzer/ebzbaybot/internal/constants"
)

type LoadOpts struct {
	Connection *sql.DB
}

func Load(opts LoadOpts) (*Stats, error) {
	connection := opts.Connection
	var chatsCount, channelsCount, collectionsCount int

	chatsCountQuery, err := connection.Query(
		"SELECT COUNT(*) FROM (SELECT DISTINCT chat_id FROM watches) AS temp;",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get chats count: %s", err)
	}
	defer chatsCountQuery.Close()
	for chatsCountQuery.Next() {
		chatsCountQuery.Scan(&chatsCount)
	}

	channelsCountQuery, err := connection.Query(
		"SELECT COUNT(*) FROM (SELECT DISTINCT chat_id FROM watches_channel) AS temp;",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get chats count: %s", err)
	}
	defer channelsCountQuery.Close()
	for channelsCountQuery.Next() {
		channelsCountQuery.Scan(&channelsCount)
	}

	collectionsCountQuery, err := connection.Query(
		"SELECT COUNT(*) FROM (SELECT DISTINCT id FROM collections) AS temp;",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get collections count: %s", err)
	}
	defer collectionsCountQuery.Close()
	for collectionsCountQuery.Next() {
		collectionsCountQuery.Scan(&collectionsCount)
	}

	goroutinesCount := runtime.NumGoroutine()
	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)
	memAlloc := memStats.Alloc / 1024 / 1024
	memTotalAlloc := memStats.TotalAlloc / 1024 / 1024
	memSystem := memStats.Sys / 1024 / 1024

	return &Stats{
		ChatsCount:        chatsCount,
		ChannelsCount:     channelsCount,
		CollectionsCount:  collectionsCount,
		GoroutinesCount:   goroutinesCount,
		AllocatedMiB:      memAlloc,
		TotalAllocatedMiB: memTotalAlloc,
		SystemMiB:         memSystem,
		Uptime:            time.Since(constants.InstanceStartedAt),
		Version:           constants.Version,
	}, nil
}
