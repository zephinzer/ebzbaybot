package handlers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/constants"
)

func HandleStart(opts Opts) error {
	msg := tgbotapi.NewMessage(opts.Update.Message.Chat.ID, fmt.Sprintf(
		"👋🏼 *HELLO/你好/HOLA/GUTEN TAG/BONJOUR*\n\n"+
			"🌟 I'm @EbzBayBot and I help you with NFT adventures at the [Ebisus Bay NFT Marketplace](https://app.ebisusbay.com) by notifying you about changes to collections you're eyeing.\n\n"+
			"⚙️ Check out my [code on GitHub](https://github.com/zephinzer/ebzbaybot)\n\n"+
			"💡 Donations address for server costs + eternal gratitude: `%s`\n\n"+
			"ℹ️ Use /help or start typing a `/` to see available commands",
		constants.DonationAddress,
	))
	msg.ParseMode = "markdown"
	msg.ReplyToMessageID = opts.Update.Message.MessageID
	_, err := opts.Bot.Send(msg)
	return err
}
