package handlers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleIDK(opts Opts) error {
	msg := tgbotapi.NewMessage(opts.Update.FromChat().ID, fmt.Sprintf(
		"👋 idk how to respond to this yet, try getting some /help?",
	))
	if opts.Update.Message != nil {
		msg.ReplyToMessageID = opts.Update.Message.MessageID
	}
	msg.ParseMode = "markdown"
	_, err := opts.Bot.Send(msg)
	return err
}
