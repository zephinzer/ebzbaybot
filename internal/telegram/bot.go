package telegram

import (
	"database/sql"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/lifecycle"
	"github.com/zephinzer/ebzbaybot/internal/telegram/handlers"
	"github.com/zephinzer/ebzbaybot/internal/utils/log"
)

type StartBotOpts struct {
	ApiKey         string
	Connection     *sql.DB
	IsDebugEnabled bool
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
		lifecycle.StartUpdatingWatchers(lifecycle.WatchingOpts{
			Bot:        bot,
			Connection: opts.Connection,
		})
	}()
	go func() {
		lifecycle.StartUpdatingChannelWatchers(lifecycle.WatchingOpts{
			Bot:        bot,
			Connection: opts.Connection,
		})
	}()
	for update := range updates {
		if update.ChannelPost != nil {
			collectionIdentifier := update.ChannelPost.CommandArguments()
			switch update.ChannelPost.Command() {
			case "init":
				if collectionIdentifier == "" {
					if _, err := bot.Send(tgbotapi.NewMessageToChannel(
						"@"+update.ChannelPost.Chat.UserName,
						"Missing a collection identifier!",
					)); err != nil {
						log.Warnf("failed to send message to channel: %s", err)
					}
				}
				handlers.HandleWatch(handlers.Opts{
					Update:     update,
					Bot:        bot,
					Connection: opts.Connection,
				})
			case "start":
				if _, err := bot.Send(tgbotapi.NewMessageToChannel(
					"@"+update.ChannelPost.Chat.UserName,
					"Initialise this channel by using /init followed by the collection identifier",
				)); err != nil {
					log.Warnf("failed to send message to channel: %s", err)
				}
			}
		}
		if update.CallbackQuery != nil {
			handleCallback(
				update,
				bot,
				opts.Connection,
			)
		}
		if update.Message == nil {
			continue
		}
		command := update.Message.Command()
		switch true {
		case command != "":
			handleCommand(update, bot, opts.Connection)
		default:
			log.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
	return nil
}
