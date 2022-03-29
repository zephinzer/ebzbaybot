package handlers

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/collection"
	"github.com/zephinzer/ebzbaybot/internal/utils/log"
	"github.com/zephinzer/ebzbaybot/pkg/ebzbay"
)

const (
	CALLBACK_LIST_GET = "list/get"
)

func handleListCallback(opts Opts) error {
	callbackData := opts.Update.CallbackQuery.Data
	callback := strings.Split(callbackData, "/")
	callbackAction := callback[1]
	switch callbackAction {
	case "get":
		chatID := opts.Update.FromChat().ID
		if len(callback) < 3 {
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
				"⚠️ Not sure why I did not find a collection address in my callback, probably ping the devs! Use /help to find out how.",
			))
			msg.ParseMode = "markdown"
			msg.ReplyToMessageID = opts.Update.CallbackQuery.Message.MessageID
			_, err := opts.Bot.Send(msg)
			return err
		}
		callbackCollection := callback[2]
		collectionDetails, _ := collection.GetCollectionByIdentifier(callbackCollection)

		collectionStats := ebzbay.GetCollectionStats(callbackCollection)

		// remove the old message
		deleteMessageRequest := tgbotapi.NewDeleteMessage(chatID, opts.Update.CallbackQuery.Message.MessageID)
		opts.Bot.Send(deleteMessageRequest)

		return sendCollectionDetails(opts, collectionStats, collectionDetails)
	}
	return nil
}

func handleListCommand(opts Opts) error {
	collectionsInlineKeyboard := getCollectionsAsKeyboard(CALLBACK_LIST_GET)
	msg := tgbotapi.NewMessage(opts.Update.Message.Chat.ID,
		"👋🏼 Here are collections I know about. If you don't see a collection you're after, it's probably not whitelisted. Use /help to find out how you can let me know about them.",
	)
	msg.ReplyMarkup = collectionsInlineKeyboard
	msg.ParseMode = "markdown"
	_, err := opts.Bot.Send(msg)
	if err != nil {
		log.Warnf("failed to send response: %s", err)
	}
	return err
}

func HandleList(opts Opts) error {
	if opts.Update.CallbackQuery != nil {
		return handleListCallback(opts)
	}
	return handleListCommand(opts)
}
