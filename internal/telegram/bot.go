package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/lifecycle"
	"github.com/zephinzer/ebzbaybot/internal/storage"
	"github.com/zephinzer/ebzbaybot/internal/utils/log"
)

type StartBotOpts struct {
	ApiKey         string
	IsDebugEnabled bool
	Storage        storage.Storage
}

func StartBot(opts StartBotOpts) error {
	bot, err := tgbotapi.NewBotAPI(opts.ApiKey)
	if err != nil {
		return fmt.Errorf("failed to instantiate a new bot: %s", err)
	}
	bot.Debug = opts.IsDebugEnabled
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := bot.GetUpdatesChan(u)
	log.Infof("starting bot at https://t.me/%s...", bot.Self.UserName)
	go func() {
		lifecycle.StartUpdatingWatchers(lifecycle.WatchingOpts{Bot: bot, Storage: opts.Storage})
	}()
	for update := range updates {
		if update.CallbackQuery != nil {
			handleCallback(update, bot, opts.Storage)
		}
		if update.Message == nil {
			continue
		}
		command := update.Message.Command()
		switch true {
		case command != "":
			handleCommand(update, bot)
		default:
			log.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
	return nil
}
