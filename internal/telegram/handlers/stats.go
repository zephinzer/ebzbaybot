package handlers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zephinzer/ebzbaybot/internal/constants"
	"github.com/zephinzer/ebzbaybot/internal/stats"
)

func HandleStats(opts Opts) error {
	statistics, err := stats.Load(stats.LoadOpts{
		Connection: opts.Connection,
	})
	if err != nil {
		return fmt.Errorf("failed to load statistics: %s", err)
	}
	msg := tgbotapi.NewMessage(opts.Update.Message.Chat.ID, fmt.Sprintf(
		"\n*🌟 Value metrics*\n"+
			"Number of Users: %v\n"+
			"Number of Channels: %v\n"+
			"Number of Collections: %v\n"+
			"\n⚙️ *System metrics*\n"+
			"Alloc (MiB): %v\n"+
			"Total (MiB): %v\n"+
			"System (MiB): %v\n"+
			"Uptime: %v\n"+
			"\n*📦 Product information*\n"+
			"Version: %s\n"+
			"\n*🙏🏼 Made by*\n"+
			"$CRO/$ETH Address: `%s`",
		statistics.ChatsCount,
		statistics.ChannelsCount,
		statistics.CollectionsCount,
		statistics.AllocatedMiB,
		statistics.TotalAllocatedMiB,
		statistics.SystemMiB,
		statistics.Uptime,
		statistics.Version,
		constants.DonationAddress,
	))
	msg.ParseMode = "markdown"
	msg.ReplyToMessageID = opts.Update.Message.MessageID
	_, err = opts.Bot.Send(msg)
	return err
}
