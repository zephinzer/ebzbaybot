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
	CALLBACK_WATCH_CONFIRM           = "watch/confirm"
	CALLBACK_WATCH_CONFIRM_NO_DELETE = "watch/confirm-retain"
	CALLBACK_WATCH_CANCEL            = "watch/cancel"
	CALLBACK_WATCH_SELECT            = "watch/select"
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
	if opts.Update.ChannelPost != nil {
		return handleWatchChannel(opts)
	}
	return handleWatchCommand(opts)
}

func handleWatchCallback(opts Opts) error {
	removeOriginalMessage := true
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
		if removeOriginalMessage {
			deleteMessageRequest := tgbotapi.NewDeleteMessage(chatID, opts.Update.CallbackQuery.Message.MessageID)
			opts.Bot.Send(deleteMessageRequest)
		}
		return handleWatchConfirmation(opts, callbackCollection)
	case "confirm-retain":
		removeOriginalMessage = false
		fallthrough
	case "confirm":
		chatID := opts.Update.FromChat().ID
		if removeOriginalMessage {
			deleteMessageRequest := tgbotapi.NewDeleteMessage(chatID, opts.Update.CallbackQuery.Message.MessageID)
			opts.Bot.Send(deleteMessageRequest)
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

		log.Infof("saving watch to db...")
		databaseWatchInstances := watch.Watches{
			watch.Watch{
				ChatID:       chatID,
				CollectionID: callbackCollection,
				LastUpdated:  time.Now(),
			},
		}
		if err := watch.Save(watch.SaveOpts{
			Connection: opts.Connection,
			Watches:    databaseWatchInstances,
		}); err != nil {
			log.Warnf("failed to save watch in chat[%v] for collection[%s]: %s", chatID, callbackCollection, err)
		}

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

func handleWatchChannel(opts Opts) error {
	chatID := opts.Update.ChannelPost.Chat.UserName
	collectionID := opts.Update.ChannelPost.CommandArguments()
	collectionDetails, err := collection.GetCollectionByIdentifier(collectionID)
	if err != nil {
		log.Warnf("failed to get collection identified by '%s': %s", collectionID, err)
		return fmt.Errorf("failed to get collection")
	}

	if err := watch.SaveChannel(watch.SaveChannelOpts{
		Connection: opts.Connection,
		Watches: watch.ChannelWatches{
			{
				ChatID:       chatID,
				CollectionID: collectionDetails.ID,
			},
		},
	}); err != nil {
		log.Warnf("failed to save channel watch on collection[%s] for channel[%s]: %s", collectionDetails.ID, chatID, err)
		_, err = opts.Bot.Send(tgbotapi.NewMessageToChannel(
			"@"+chatID,
			fmt.Sprintf("Something went wrong while trying to watch the *%s* collection. See logs for more information", collectionDetails.Label),
		))
		return err
	}
	msg := tgbotapi.NewMessageToChannel(
		"@"+chatID,
		fmt.Sprintf("I will notify this channel for updates on the *%s* collection", collectionDetails.Label),
	)
	msg.ParseMode = "markdown"
	_, err = opts.Bot.Send(msg)
	return err
}

func handleWatchCommand(opts Opts) error {
	collectionIdentifier := opts.Update.Message.CommandArguments()
	if collectionIdentifier == "" {
		collectionsInlineKeyboard, err := getUnwatchedCollectionsAsKeyboard(CALLBACK_WATCH_SELECT, opts.Update.FromChat().ID, opts.Connection)
		if err != nil {
			return fmt.Errorf("failed to get collections as a keyboard: %s", err)
		}
		msg := tgbotapi.NewMessage(opts.Update.FromChat().ID, fmt.Sprintf(
			"👋🏼 I sense interest in you, young one. Which collection shall I update you about?",
		))
		msg.ParseMode = "markdown"
		msg.ReplyToMessageID = opts.Update.Message.MessageID
		msg.ReplyMarkup = collectionsInlineKeyboard
		_, err = opts.Bot.Send(msg)
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
		collectionDetails.Label,
		collectionDetails.ID,
		collectionDetails.ID,
		collectionDetails.ID,
	))
	msg.ParseMode = "markdown"
	if opts.Update.Message != nil {
		msg.ReplyToMessageID = opts.Update.Message.MessageID
	}
	msg.ReplyMarkup = getWatchKeyboard(collectionDetails.ID)
	_, err = opts.Bot.Send(msg)
	return err
}
