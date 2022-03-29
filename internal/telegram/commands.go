package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/telegram/handlers"
	"github.com/zephinzer/ebzbaybot/internal/utils/log"
)

func handleCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	log.Infof("chat[%v]: command[%s] %s", update.Message.Chat.ID, update.Message.Command(), update.Message.CommandArguments())
	opts := handlers.Opts{Update: update, Bot: bot}
	switch update.Message.Command() {
	case "start":
		return handlers.HandleStart(opts)
	case "list":
		return handlers.HandleList(opts)
	case "get":
		return handlers.HandleGet(opts)
	case "watch":
		return handlers.HandleWatch(opts)
	case "help":
		return handlers.HandleHelp(opts)
	}
	return handlers.HandleIDK(opts)
}
