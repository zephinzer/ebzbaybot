package lifecycle

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/storage"
)

type ScrapingOpts struct {
	ScrapeInterval time.Duration
	Storage        storage.Storage
}

type WatchingOpts struct {
	Bot     *tgbotapi.BotAPI
	Storage storage.Storage
}
