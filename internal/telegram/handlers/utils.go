package handlers

import (
	"database/sql"
	"fmt"
	"path"
	"sort"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/collection"
	"github.com/zephinzer/ebzbaybot/internal/utils/log"
	"github.com/zephinzer/ebzbaybot/internal/watch"
	"github.com/zephinzer/ebzbaybot/pkg/constants"
	"github.com/zephinzer/ebzbaybot/pkg/ebzbay"
)

func getWatchedCollectionsAsKeyboard(callbackActionPrefix string, chatID int64, connection *sql.DB) (tgbotapi.InlineKeyboardMarkup, error) {
	watches, err := watch.Load(watch.LoadOpts{
		Connection: connection,
		OnlyFor:    chatID,
	})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("failed to load watches for chat[%v]: %s", chatID, err)
	}
	watchesMap := map[string]bool{}
	for _, watchInstance := range watches {
		watchesMap[watchInstance.CollectionID] = true
	}

	inlineKeyboardButtons := [][]tgbotapi.InlineKeyboardButton{}
	collectionNames := []string{}
	for address, names := range constants.CollectionByAddress {
		if isWatching, exists := watchesMap[address]; isWatching && exists {
			collectionNames = append(collectionNames, names[0])
		}
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
	return tgbotapi.NewInlineKeyboardMarkup(inlineKeyboardButtons...), nil
}

func getUnwatchedCollectionsAsKeyboard(callbackActionPrefix string, chatID int64, connection *sql.DB) (tgbotapi.InlineKeyboardMarkup, error) {
	watches, err := watch.Load(watch.LoadOpts{
		Connection: connection,
		OnlyFor:    chatID,
	})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("failed to load watches for chat[%v]: %s", chatID, err)
	}
	watchesMap := map[string]bool{}
	for _, watchInstance := range watches {
		watchesMap[watchInstance.CollectionID] = true
	}

	inlineKeyboardButtons := [][]tgbotapi.InlineKeyboardButton{}
	collectionNames := []string{}
	for address, names := range constants.CollectionByAddress {
		if isWatching, exists := watchesMap[address]; !isWatching && !exists {
			collectionNames = append(collectionNames, names[0])
		}
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
	return tgbotapi.NewInlineKeyboardMarkup(inlineKeyboardButtons...), nil
}

func getCollectionsAsKeyboard(callbackActionPrefix string) (tgbotapi.InlineKeyboardMarkup, error) {
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
	return tgbotapi.NewInlineKeyboardMarkup(inlineKeyboardButtons...), nil
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
	if !opts.Update.FromChat().IsChannel() {
		watchExists, err := watch.Exists(watch.ExistsOpts{
			Connection:   opts.Connection,
			ChatID:       chatID,
			CollectionID: details.ID,
		})
		if err != nil {
			log.Warnf("failed to check whether watch exists for chat[%s] and collection[%s]: %s", chatID, details.ID, err)
		}
		if !watchExists {
			responseMessage.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(
						"WATCH THIS COLLECTION",
						path.Join(
							CALLBACK_WATCH_CONFIRM_NO_DELETE,
							details.ID,
						),
					),
				),
			)
		}
	}
	_, err := opts.Bot.Send(responseMessage)
	return err
}
