package handlers

import (
	"encoding/json"
	"fmt"
	"path"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/collection"
	"github.com/zephinzer/ebzbaybot/internal/types"
	"github.com/zephinzer/ebzbaybot/pkg/constants"
)

const (
	CALLBACK_WATCH_CONFIRM = "watch/confirm"
	CALLBACK_WATCH_CANCEL  = "watch/cancel"
	CALLBACK_WATCH_SELECT  = "watch/select"
)

func getWatchKeyboard(collection string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("WATCH", path.Join(CALLBACK_WATCH_CONFIRM, collection)),
			tgbotapi.NewInlineKeyboardButtonData("CANCEL", CALLBACK_WATCH_CANCEL),
		),
	)
}

func HandleWatch(opts Opts) error {
	if opts.Update.CallbackQuery != nil {
		return handleWatchCallback(opts)
	}
	return handleWatchCommand(opts)
}

func handleWatchCallback(opts Opts) error {
	callbackData := opts.Update.CallbackQuery.Data
	callback := strings.Split(callbackData, "/")
	callbackAction := callback[1]
	switch callbackAction {
	case "select":
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
		deleteMessageRequest := tgbotapi.NewDeleteMessage(chatID, opts.Update.CallbackQuery.Message.MessageID)
		opts.Bot.Send(deleteMessageRequest)
		return handleWatchConfirmation(opts, callbackCollection)
	case "confirm":
		chatID := opts.Update.FromChat().ID
		deleteMessageRequest := tgbotapi.NewDeleteMessage(chatID, opts.Update.CallbackQuery.Message.MessageID)
		opts.Bot.Send(deleteMessageRequest)
		if opts.Storage == nil {
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
				"⚠️ I couldn't store this information because the dev forgot to add a storage component to me",
			))
			msg.ParseMode = "markdown"
			msg.ReplyToMessageID = opts.Update.CallbackQuery.Message.MessageID
			_, err := opts.Bot.Send(msg)
			return err
		}
		if len(callback) < 3 {
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
				"⚠️ Something messed up, sorry - I did not find a collection to watch in my callback",
			))
			msg.ParseMode = "markdown"
			msg.ReplyToMessageID = opts.Update.CallbackQuery.Message.MessageID
			_, err := opts.Bot.Send(msg)
			return err
		}

		callbackCollection := callback[2]

		userID := strconv.FormatInt(chatID, 10)
		watches := types.WatchStorage{}
		watchesJSON, _ := opts.Storage.Get("watches")
		if watchesJSON != nil {
			json.Unmarshal([]byte(watchesJSON), &watches)
		}
		if _, exist := watches[userID]; !exist {
			watches[userID] = types.Watch{
				CollectionMap: map[string]string{},
			}
		}
		currentTimestamp := strconv.FormatInt(time.Now().Unix(), 10)
		watches[userID].CollectionMap[callbackCollection] = currentTimestamp
		watchesJSON, _ = json.Marshal(watches)
		opts.Storage.Set("watches", watchesJSON)

		collectionName := constants.CollectionByAddress[callbackCollection][0]
		msg := tgbotapi.NewMessage(opts.Update.FromChat().ID, fmt.Sprintf(
			"😍 I will notify you on changes to the [%s collection](https://app.ebisusbay.com/collection/%s) from now on! You may use /unwatch to unregister this should your interest in *%s* wane",
			collectionName,
			callbackCollection,
			collectionName,
		))
		msg.ParseMode = "markdown"
		_, err := opts.Bot.Send(msg)
		return err
	case "cancel":
		deleteMessageRequest := tgbotapi.NewDeleteMessage(opts.Update.FromChat().ID, opts.Update.CallbackQuery.Message.MessageID)
		_, err := opts.Bot.Send(deleteMessageRequest)
		return err
	}
	return nil
}

func handleWatchCommand(opts Opts) error {
	collectionIdentifier := opts.Update.Message.CommandArguments()
	if collectionIdentifier == "" {
		collectionsInlineKeyboard := getCollectionsAsKeyboard(CALLBACK_WATCH_SELECT)
		msg := tgbotapi.NewMessage(opts.Update.FromChat().ID, fmt.Sprintf(
			"👋🏼 I sense interest in you, young one. Which collection shall I update you about?",
		))
		msg.ParseMode = "markdown"
		msg.ReplyToMessageID = opts.Update.Message.MessageID
		msg.ReplyMarkup = collectionsInlineKeyboard
		_, err := opts.Bot.Send(msg)
		return err
	}
	return handleWatchConfirmation(opts, collectionIdentifier)
}

func handleWatchConfirmation(opts Opts, collectionIdentifier string) error {
	collectionDetails, err := collection.GetCollectionByIdentifier(collectionIdentifier)
	if err != nil {
		msg := tgbotapi.NewMessage(opts.Update.FromChat().ID, fmt.Sprintf(
			"⚠️ Apologies, no collection could be found with the provided identifier `%s`",
			collectionIdentifier,
		))
		msg.ParseMode = "markdown"
		msg.ReplyToMessageID = opts.Update.Message.MessageID
		_, err := opts.Bot.Send(msg)
		return err
	}

	msg := tgbotapi.NewMessage(opts.Update.FromChat().ID, fmt.Sprintf(
		"👀 Interesting choice you have made. Shall I confirm *%s* collection with the token address `%s` as your destiny for today?\n\n"+
			"👉🏼 Review on [Cronoscan](https://cronoscan.com/address/%s) | [Ebisus Bay](https://app.ebisusbay.com/collection/%s)",
		collectionDetails.Name,
		collectionDetails.Address,
		collectionDetails.Address,
		collectionDetails.Address,
	))
	msg.ParseMode = "markdown"
	if opts.Update.Message != nil {
		msg.ReplyToMessageID = opts.Update.Message.MessageID
	}
	msg.ReplyMarkup = getWatchKeyboard(collectionDetails.Address)
	_, err = opts.Bot.Send(msg)
	return err
}
