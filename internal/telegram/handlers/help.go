package handlers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/constants"
)

func HandleHelp(opts Opts) error {
	msg := tgbotapi.NewMessage(opts.Update.Message.Chat.ID, fmt.Sprintf(
		"🎣 I'm @EbzBayBot and I help you with NFT adventures at the [Ebisus Bay NFT Marketplace](https://app.ebisusbay.com) by notifying you about changes to collections you're eyeing.\n\n"+
			"💡 Donations address for server costs + eternal gratitude ($CRO/$ETH): `%s`\n\n"+
			"*🪓 To control me, use*:\n"+
			"`/get` gets information about a collection\n"+
			"`/watch` starts watching a collection\n"+
			"`/unwatch` unwatches a collection\n"+
			"`/list` lists available collections\n",
		constants.DonationAddress,
	))
	msg.ParseMode = "markdown"
	_, err := opts.Bot.Send(msg)
	return err
}
