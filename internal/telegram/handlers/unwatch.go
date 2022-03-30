package handlers

import (
	"fmt"
	"path"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/collection"
	"github.com/zephinzer/ebzbaybot/internal/utils/log"
	"github.com/zephinzer/ebzbaybot/internal/watch"
	"github.com/zephinzer/ebzbaybot/pkg/constants"
)

const (
	CALLBACK_UNWATCH_CONFIRM = "unwatch/confirm"
	CALLBACK_UNWATCH_CANCEL  = "unwatch/cancel"
	CALLBACK_UNWATCH_SELECT  = "unwatch/select"
)

func getUnwatchKeyboard(collection string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("UNWATCH", path.Join(CALLBACK_UNWATCH_CONFIRM, collection)),
			tgbotapi.NewInlineKeyboardButtonData("CANCEL", CALLBACK_UNWATCH_CANCEL),
		),
	)
}

func HandleUnwatch(opts Opts) error {
	if opts.Update.CallbackQuery != nil {
		return handleUnwatchCallback(opts)
	}
	return handleUnwatchCommand(opts)
}

func handleUnwatchCallback(opts Opts) error {
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
		return handleUnwatchConfirmation(opts, callbackCollection)
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

		log.Infof("removing watch to db...")
		databaseWatchInstances := watch.Watches{
			watch.Watch{
				ChatID:       chatID,
				CollectionID: callbackCollection,
				LastUpdated:  time.Now(),
			},
		}
		watch.Remove(watch.RemoveOpts{
			Connection: opts.Connection,
			Watches:    databaseWatchInstances,
		})

		collectionName := constants.CollectionByAddress[callbackCollection][0]
		msg := tgbotapi.NewMessage(opts.Update.FromChat().ID, fmt.Sprintf(
			"💃🏻 I'll no longer be notifying you about *%s*!",
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

func handleUnwatchCommand(opts Opts) error {
	collectionIdentifier := opts.Update.Message.CommandArguments()
	if collectionIdentifier == "" {
		collectionsInlineKeyboard, err := getWatchedCollectionsAsKeyboard(
			CALLBACK_UNWATCH_SELECT,
			opts.Update.FromChat().ID,
			opts.Connection,
		)
		if err != nil {
			return fmt.Errorf("failed to get a keyboard with the watched collections: %s", err)
		}
		msg := tgbotapi.NewMessage(opts.Update.FromChat().ID, fmt.Sprintf(
			"👋🏼 Which collection would you like to unwatch?",
		))
		msg.ParseMode = "markdown"
		msg.ReplyToMessageID = opts.Update.Message.MessageID
		msg.ReplyMarkup = collectionsInlineKeyboard
		_, err = opts.Bot.Send(msg)
		return err
	}
	return handleUnwatchConfirmation(opts, collectionIdentifier)
}

func handleUnwatchConfirmation(opts Opts, collectionIdentifier string) error {
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
		"👀 Are you sure you would like to stop watching the *%s* collection?",
		collectionDetails.Label,
	))
	msg.ParseMode = "markdown"
	if opts.Update.Message != nil {
		msg.ReplyToMessageID = opts.Update.Message.MessageID
	}
	msg.ReplyMarkup = getUnwatchKeyboard(collectionDetails.ID)
	_, err = opts.Bot.Send(msg)
	return err
}
