package handlers

import (
	"database/sql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/storage"
)

type Opts struct {
	Bot        *tgbotapi.BotAPI
	Connection *sql.DB
	Storage    storage.Storage
	Update     tgbotapi.Update
}
