package handlers

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleHelp(opts Opts) error {
	var listOfCommands strings.Builder
	listOfCommands.WriteString(fmt.Sprintf("/start starts this bot\n"))
	listOfCommands.WriteString(fmt.Sprintf("/get gets information about a collection\n"))
	listOfCommands.WriteString(fmt.Sprintf("/watch lets you watch a collection\n"))
	listOfCommands.WriteString(fmt.Sprintf("/list lists available collections\n"))
	msg := tgbotapi.NewMessage(opts.Update.Message.Chat.ID, listOfCommands.String())
	msg.ParseMode = "markdown"
	_, err := opts.Bot.Send(msg)
	return err
}
