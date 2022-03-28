package telegram

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/utils/log"
)

type StartBotOpts struct {
	ApiKey         string
	IsDebugEnabled bool
}

func StartBot(opts StartBotOpts) error {
	bot, err := tgbotapi.NewBotAPI(opts.ApiKey)
	if err != nil {
		return fmt.Errorf("failed to instantiate a new bot: %s", err)
	}
	bot.Debug = opts.IsDebugEnabled

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	log.Infof("starting bot at https://t.me/%s...", bot.Self.UserName)

	lifecycleInterval := 10 * time.Second
	go func(tick <-chan time.Time) {
		for {
			<-tick

		}
	}(time.NewTicker(lifecycleInterval).C)

	for update := range updates {
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
