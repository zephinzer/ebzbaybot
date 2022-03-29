package handlers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleStart(opts Opts) error {
	msg := tgbotapi.NewMessage(opts.Update.Message.Chat.ID, fmt.Sprintf(
		"👋 there, I am a bot that helps you with your NFT adventures at the [Ebisus Bay NFT marketplace](https://app.ebisusbay.com).\n\n"+
			"I am an open-source software and my source code can be found at https://github.com/zephinzer/ebzbaybot, feel free to raise a pull request with any feature you want.\n\n"+
			"Start typing a command to see available commands or use /help to see all commands in a single message",
	))
	msg.ParseMode = "markdown"
	msg.ReplyToMessageID = opts.Update.Message.MessageID
	_, err := opts.Bot.Send(msg)
	return err
}
