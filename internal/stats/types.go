package stats

import "time"

type Stats struct {
	// service related

	ChannelsCount    int `json:"channelsCount"`
	ChatsCount       int `json:"chatsCount"`
	CollectionsCount int `json:"collectionsCount"`
	GoroutinesCount  int `json:"goroutinesCount"`

	// system related

	AllocatedMiB      uint64        `json:"allocatedMemoryMiB"`
	TotalAllocatedMiB uint64        `json:"totalAllocatedMemoryMB"`
	SystemMiB         uint64        `json:"systemMemoryMiB"`
	Uptime            time.Duration `json:"uptime"`

	// code related

	Version string `json:"version"`
}
