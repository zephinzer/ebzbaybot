package handlers

import (
	"fmt"
	"path"
	"sort"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/collection"
	"github.com/zephinzer/ebzbaybot/pkg/constants"
	"github.com/zephinzer/ebzbaybot/pkg/ebzbay"
)

func getCollectionsAsKeyboard(callbackActionPrefix string) tgbotapi.InlineKeyboardMarkup {
	inlineKeyboardButtons := [][]tgbotapi.InlineKeyboardButton{}
	collectionNames := []string{}
	for _, names := range constants.CollectionByAddress {
		collectionNames = append(collectionNames, names[0])
	}
	sort.Strings(collectionNames)
	for _, collectionName := range collectionNames {
		collectionInstance, _ := collection.GetCollectionByIdentifier(collectionName)
		inlineKeyboardButtons = append(
			inlineKeyboardButtons,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(collectionName, path.Join(callbackActionPrefix, collectionInstance.ID)),
			),
		)
	}
	return tgbotapi.NewInlineKeyboardMarkup(inlineKeyboardButtons...)
}

func sendCollectionDetails(opts Opts, stats *ebzbay.CollectionStats, details *collection.Collection) error {
	chatID := opts.Update.FromChat().ID

	aliasText := ""
	aliases := constants.CollectionByAddress[details.ID]
	if len(aliases) > 1 {
		aliasText = fmt.Sprintf("👤 Aliases: `%s`\n", strings.Join(aliases[1:], "`, `"))
	}

	responseMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf(
		"🎨 *%s* | Token address: `%s`\n"+
			"%s"+
			"📊 Total listings: %v\n"+
			"💰 Floor price: *%v* $CRO\n"+
			"⚓️ Average 10 lowest prices: _%v $CRO_\n\n"+
			"👉🏼 View on [Cronoscan](https://cronoscan.com/address/%s) | [Ebisus Bay](https://app.ebisusbay.com/collection/%s)",
		details.Label,
		details.ID,
		aliasText,
		stats.Listings,
		stats.FloorPrice,
		stats.AverageLowestTenPrice,
		details.ID,
		details.ID,
	))
	responseMessage.ParseMode = "markdown"
	_, err := opts.Bot.Send(responseMessage)
	return err
}
