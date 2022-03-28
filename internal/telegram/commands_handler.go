package telegram

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/utils/log"
	"github.com/zephinzer/ebzbaybot/pkg/constants"
	"github.com/zephinzer/ebzbaybot/pkg/ebzbay"
)

func handleCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	log.Infof("chat[%v]: command[%s] %s", update.Message.Chat.ID, update.Message.Command(), update.Message.CommandArguments())
	switch update.Message.Command() {
	case "start":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(
			"👋 there, I am a bot that helps you with your NFT adventures at the [Ebisus Bay NFT marketplace](https://app.ebisusbay.com).\n\n"+
				"I am an open-source software and my source code can be found at https://github.com/zephinzer/ebzbaybot, feel free to raise a pull request with any feature you want.\n\n"+
				"Start typing a command to see available commands or use /help to see all commands in a single message",
		))
		msg.ParseMode = "markdown"
		msg.ReplyToMessageID = update.Message.MessageID
		_, err := bot.Send(msg)
		return err
	case "list":
		var listOfCollections strings.Builder
		for address, name := range constants.Collection {
			listOfCollections.WriteString(fmt.Sprintf("%s: `/get %s`\n", name, address))
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, listOfCollections.String())
		msg.ParseMode = "markdown"
		_, err := bot.Send(msg)
		return err
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
			_, err := bot.Send(msg)
			return err
		}

		collectionName := constants.Collection[collectionAddress]
		if collectionName == "" {
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				fmt.Sprintf(
					"⚠️ The provided address is not whitelisted, use /list to see available collections.\n\n"+
						"If you would like to add this collection, you can do so by raising a pull request adding your collection into [this file](https://github.com/zephinzer/ebzbaybot/blob/master/pkg/constants/data.json)",
				),
			)
			msg.ParseMode = "markdown"
			msg.ReplyToMessageID = update.Message.MessageID
			_, err := bot.Send(msg)
			return err
		}

		waitingMessage := tgbotapi.NewChatAction(update.Message.Chat.ID, "typing")
		if _, err := bot.Send(waitingMessage); err != nil {
			log.Warnf("failed to send in-progress action: %s", err)
		}

		collectionStats := ebzbay.GetCollectionStats(update.Message.CommandArguments())
		responseMessage := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(
			"Collection name: **%s**\n"+
				"Collection address: `%s`\n"+
				"Total listings: %v\n"+
				"Floor price: %v CRO\n"+
				"Average 10 lowest prices: %v CRO\n",
			collectionName,
			collectionAddress,
			collectionStats.Listings,
			collectionStats.FloorPrice,
			collectionStats.AverageLowestTenPrice,
		))
		responseMessage.ParseMode = "markdown"
		responseMessage.ReplyToMessageID = update.Message.MessageID
		_, err := bot.Send(responseMessage)
		return err
	case "help":
		var listOfCommands strings.Builder
		listOfCommands.WriteString(fmt.Sprintf("/start starts this bot\n"))
		listOfCommands.WriteString(fmt.Sprintf("/get gets information about a collection\n"))
		listOfCommands.WriteString(fmt.Sprintf("/list lists available collections\n"))
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, listOfCommands.String())
		msg.ParseMode = "markdown"
		bot.Send(msg)
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(
			"👋 idk how to respond to this, try getting some /help",
		))
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ParseMode = "markdown"
		_, err := bot.Send(msg)
		return err
	}
	return nil
}
