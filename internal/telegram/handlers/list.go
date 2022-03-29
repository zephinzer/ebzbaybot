package handlers

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/pkg/constants"
)

func HandleList(opts Opts) error {
	var listOfCollections strings.Builder
	for address, names := range constants.CollectionByAddress {
		listOfCollections.WriteString(fmt.Sprintf("%s: `/get %s`\n", names[0], address))
	}
	msg := tgbotapi.NewMessage(opts.Update.Message.Chat.ID, listOfCollections.String())
	msg.ParseMode = "markdown"
	_, err := opts.Bot.Send(msg)
	return err
}
