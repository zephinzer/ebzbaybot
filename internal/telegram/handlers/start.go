package handlers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleStart(opts Opts) error {
	msg := tgbotapi.NewMessage(opts.Update.Message.Chat.ID, fmt.Sprintf(
		"👋 there, I am @EbzBayBot, a bot that helps you with your NFT adventures at the [Ebisus Bay NFT marketplace](https://app.ebisusbay.com).\n\n"+
			"My creator has open-sourced me and my source code can be found at https://github.com/zephinzer/ebzbaybot, feel free to raise a pull request with features/collection whitelist requests.\n\n"+
			"Use /help to see what I can help you with in a single message",
	))
	msg.ParseMode = "markdown"
	msg.ReplyToMessageID = opts.Update.Message.MessageID
	_, err := opts.Bot.Send(msg)
	return err
}
