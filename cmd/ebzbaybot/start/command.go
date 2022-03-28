package start

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/usvc/go-config"
	"github.com/zephinzer/ebzbaybot/internal/telegram"
)

const (
	ConfigKeyTelegramBotApiKey = "telegram-bot-api-key"
)

var conf = config.Map{
	ConfigKeyTelegramBotApiKey: &config.String{
		Shorthand: "k",
		Usage:     "API key for Telegram (get it from: https://t.me/BotFather)",
	},
}

func GetCommand() *cobra.Command {
	command := cobra.Command{
		Use:   "start",
		Short: "Starts the ebisus's bay bot",
		Long:  "Starts the ebisus's bay bot",
		RunE:  runE,
	}
	conf.ApplyToCobra(&command)
	return &command
}

func runE(cmd *cobra.Command, args []string) error {
	botToken := conf.GetString(ConfigKeyTelegramBotApiKey)
	if botToken == "" {
		return fmt.Errorf("failed to receive a telegram bot token, pass it in with --%s", ConfigKeyTelegramBotApiKey)
	}
	telegram.StartBot(telegram.StartBotOpts{
		ApiKey:         botToken,
		IsDebugEnabled: false,
	})
	return nil
}
