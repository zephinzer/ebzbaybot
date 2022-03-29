package handlers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/collection"
	"github.com/zephinzer/ebzbaybot/internal/utils/log"
	"github.com/zephinzer/ebzbaybot/pkg/ebzbay"
)

func HandleGet(opts Opts) error {
	collectionIdentifier := opts.Update.Message.CommandArguments()
	if collectionIdentifier == "" {
		msg := tgbotapi.NewMessage(
			opts.Update.Message.Chat.ID,
			fmt.Sprintf(
				"⚠️ Use the collection address or collection shortcut as an argument to this command\n\n"+
					"Example for Mad Meerkats: `/get 0x89dbc8bd9a6037cbd6ec66c4bf4189c9747b1c56`\n"+
					"Example for MM Treehouse: `/get mmt`\n",
			),
		)
		msg.ParseMode = "markdown"
		msg.ReplyToMessageID = opts.Update.Message.MessageID
		_, err := opts.Bot.Send(msg)
		return err
	}

	collectionDetails, err := collection.GetCollectionByIdentifier(collectionIdentifier)
	if err != nil {
		msg := tgbotapi.NewMessage(
			opts.Update.Message.Chat.ID,
			fmt.Sprintf(
				"⚠️ The provided identifier ('%s') is not whitelisted, use /list to see available collections.\n\n"+
					"If you would like to add this collection, you can do so by raising a pull request adding your collection into [this file](https://github.com/zephinzer/ebzbaybot/blob/master/pkg/constants/data.json)",
				collectionIdentifier,
			),
		)
		msg.ParseMode = "markdown"
		msg.ReplyToMessageID = opts.Update.Message.MessageID
		_, err := opts.Bot.Send(msg)
		return err
	}

	waitingMessage := tgbotapi.NewChatAction(opts.Update.Message.Chat.ID, "typing")
	if _, err := opts.Bot.Send(waitingMessage); err != nil {
		log.Warnf("failed to send in-progress action: %s", err)
	}

	collectionStats := ebzbay.GetCollectionStats(opts.Update.Message.CommandArguments())
	responseMessage := tgbotapi.NewMessage(opts.Update.Message.Chat.ID, fmt.Sprintf(
		"Collection name: **%s**\n"+
			"Collection address: `%s`\n"+
			"Total listings: %v\n"+
			"Floor price: %v CRO\n"+
			"Average 10 lowest prices: %v CRO\n",
		collectionDetails.Name,
		collectionDetails.Address,
		collectionStats.Listings,
		collectionStats.FloorPrice,
		collectionStats.AverageLowestTenPrice,
	))
	responseMessage.ParseMode = "markdown"
	responseMessage.ReplyToMessageID = opts.Update.Message.MessageID
	_, err = opts.Bot.Send(responseMessage)
	return err
}
