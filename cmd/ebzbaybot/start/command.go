package start

import (
	"fmt"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/usvc/go-config"
	"github.com/zephinzer/ebzbaybot/internal/lifecycle"
	"github.com/zephinzer/ebzbaybot/internal/storage"
	"github.com/zephinzer/ebzbaybot/internal/telegram"
	"github.com/zephinzer/ebzbaybot/internal/utils/log"
)

const (
	ConfigKeyTelegramBotApiKey = "telegram-bot-api-key"
	ConfigKeyStorageConfig     = "storage-config"
)

var conf = config.Map{
	ConfigKeyTelegramBotApiKey: &config.String{
		Shorthand: "k",
		Usage:     "API key for Telegram (get it from: https://t.me/BotFather)",
	},
	ConfigKeyStorageConfig: &config.String{
		Shorthand: "s",
		Default:   "memory",
		Usage:     "Define the storage type",
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
	var storageInstance storage.Storage
	switch conf.GetString(ConfigKeyStorageConfig) {
	case "memory":
		fallthrough
	default:
		storageInstance = storage.NewMemory()
	}

	var waiter sync.WaitGroup

	waiter.Add(1)

	go func() {
		lifecycleInterval := 10 * time.Second
		lifecycle.StartCollectionsScraping(lifecycle.ScrapingOpts{
			ScrapeInterval: lifecycleInterval,
			Storage:        storageInstance,
		})
		waiter.Done()
	}()

	go func() {
		if err := telegram.StartBot(telegram.StartBotOpts{
			ApiKey:         botToken,
			IsDebugEnabled: false,
			Storage:        storageInstance,
		}); err != nil {
			log.Errorf("failed to keep bot alive: %s", err)
			waiter.Done()
		}
	}()

	waiter.Wait()
	return nil
}
