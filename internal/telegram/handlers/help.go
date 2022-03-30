package handlers

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/constants"
)

func HandleHelp(opts Opts) error {
	var listOfCommands strings.Builder
	listOfCommands.WriteString("🎣 @EbzBayBot is an unofficial Telegram bot for getting notifications on floor price changes from [Ebisus Bay Marketplace](https://app.ebisusbay.com).\n")
	listOfCommands.WriteString("👨‍💻 Bot is open-sourced and available at https://github.com/zephinzer/ebzbaybot if anyone wants to contribute.\n")
	listOfCommands.WriteString("💡 Collections are available on a whitelist basis to keep everyone safe, the repository's `README.md` has instructions on how to get a collection whitelisted.")
	listOfCommands.WriteString(fmt.Sprintf("🙇 Servers cost money to run; if this bot helps you and you are financially able to, $CRO donations will be very appreciated at `%s` 🙇\n\n", constants.DonationAddress))
	listOfCommands.WriteString("🪓 Available commands are:\n")
	listOfCommands.WriteString(fmt.Sprintf("`/start` starts this bot\n"))
	listOfCommands.WriteString(fmt.Sprintf("`/get` gets information about a collection\n"))
	listOfCommands.WriteString(fmt.Sprintf("`/watch` starts watching a collection\n"))
	listOfCommands.WriteString(fmt.Sprintf("`/unwatch` unwatches a collection\n"))
	listOfCommands.WriteString(fmt.Sprintf("`/list` lists available collections\n"))
	msg := tgbotapi.NewMessage(opts.Update.Message.Chat.ID, listOfCommands.String())
	msg.ParseMode = "markdown"
	_, err := opts.Bot.Send(msg)
	return err
}
