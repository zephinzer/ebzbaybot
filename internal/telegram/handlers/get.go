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
		return HandleList(opts)
	}

	collectionDetails, err := collection.GetCollectionByIdentifier(collectionIdentifier)
	if err != nil {
		msg := tgbotapi.NewMessage(
			opts.Update.FromChat().ID,
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

	waitingMessage := tgbotapi.NewChatAction(opts.Update.FromChat().ID, "typing")
	if _, err := opts.Bot.Send(waitingMessage); err != nil {
		log.Warnf("failed to send in-progress action: %s", err)
	}

	collectionStats := ebzbay.GetCollectionStats(collectionDetails.ID)
	return sendCollectionDetails(opts, collectionStats, collectionDetails)
}
