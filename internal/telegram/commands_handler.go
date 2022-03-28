package telegram

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/utils/log"
	"github.com/zephinzer/ebzbaybot/pkg/constants"
	"github.com/zephinzer/ebzbaybot/pkg/ebzbay"
)

func handleCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	switch update.Message.Command() {
	case "start":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID
	case "list":
		var listOfCollections strings.Builder
		for address, name := range constants.Collection {
			listOfCollections.WriteString(fmt.Sprintf("%s: `/get %s`\n", name, address))
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, listOfCollections.String())
		msg.ParseMode = "markdown"
		bot.Send(msg)
	case "get":
		collectionAddress := update.Message.CommandArguments()
		if collectionAddress == "" {
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				fmt.Sprintf(
					"⚠️ Use the collection address as an argument to this command\n"+
						"Example for Mad Meerkats: `/get 0x89dbc8bd9a6037cbd6ec66c4bf4189c9747b1c56`",
				),
			)
			msg.ParseMode = "markdown"
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}

		collectionName := constants.Collection[collectionAddress]
		if collectionName == "" {
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				fmt.Sprintf(
					"⚠️ The provided address is not whitelisted, use /list to see available collections.\n\n" +
					"If you would like to add this collection, you can do so by raising a merge request at https://github.com/zephinzer/ebzbot"
				)
			)
		}

		waitingMessage := tgbotapi.NewChatAction(update.Message.Chat.ID, "typing")
		if _, err := bot.Send(waitingMessage); err != nil {
			log.Warnf("failed to send in-progress action: %s", err)
		}

		collectionStats := ebzbay.GetCollectionStats(update.Message.CommandArguments())
		responseMessage := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(
			"Total listings: %v\n"+
				"Floor price: %v CRO\n"+
				"Average 10 lowest prices: %v CRO\n",
			collectionStats.Listings,
			collectionStats.FloorPrice,
			collectionStats.AverageLowestTenPrice,
		))
		responseMessage.ReplyToMessageID = update.Message.MessageID
		bot.Send(responseMessage)
	case "":
	}
	log.Infof("chat[%v]: command[%s] %s", update.Message.Chat.ID, update.Message.Command(), update.Message.CommandArguments())
}
