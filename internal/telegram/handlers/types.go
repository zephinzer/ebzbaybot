package handlers

import (
	"database/sql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Opts struct {
	Bot        *tgbotapi.BotAPI
	Connection *sql.DB
	Update     tgbotapi.Update
}
