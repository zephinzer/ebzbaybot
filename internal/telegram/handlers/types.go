package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/storage"
)

type Opts struct {
	Bot     *tgbotapi.BotAPI
	Storage storage.Storage
	Update  tgbotapi.Update
}
