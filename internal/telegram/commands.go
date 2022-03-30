package telegram

import (
	"database/sql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/telegram/handlers"
	"github.com/zephinzer/ebzbaybot/internal/utils/log"
)

func handleCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI, connection *sql.DB) error {
	log.Infof("chat[%v]: command[%s] %s", update.Message.Chat.ID, update.Message.Command(), update.Message.CommandArguments())
	opts := handlers.Opts{
		Bot:        bot,
		Connection: connection,
		Update:     update,
	}
	var err error
	switch update.Message.Command() {
	case "start":
		err = handlers.HandleStart(opts)
	case "list":
		err = handlers.HandleList(opts)
	case "get":
		err = handlers.HandleGet(opts)
	case "watch":
		err = handlers.HandleWatch(opts)
	case "unwatch":
		err = handlers.HandleUnwatch(opts)
	case "help":
		err = handlers.HandleHelp(opts)
	default:
		err = handlers.HandleIDK(opts)
	}
	if err != nil {
		log.Warnf("a handler failed: %s", err)
	}
	return err
}
