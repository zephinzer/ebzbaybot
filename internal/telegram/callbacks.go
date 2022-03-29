package telegram

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/storage"
	"github.com/zephinzer/ebzbaybot/internal/telegram/handlers"
	"github.com/zephinzer/ebzbaybot/internal/utils/log"
)

func handleCallback(update tgbotapi.Update, bot *tgbotapi.BotAPI, storageInstance storage.Storage) error {
	log.Infof("callback[%s] %s", update.CallbackQuery.ID, update.CallbackQuery.Data)
	opts := handlers.Opts{Update: update, Bot: bot, Storage: storageInstance}
	callbackData := strings.Split(update.CallbackQuery.Data, "/")
	if len(callbackData) == 0 {
		return handlers.HandleIDK(opts)
	}
	handler := callbackData[0]
	switch handler {
	case "list":
		return handlers.HandleList(opts)
	case "watch":
		return handlers.HandleWatch(opts)
	}
	return handlers.HandleIDK(opts)
}
