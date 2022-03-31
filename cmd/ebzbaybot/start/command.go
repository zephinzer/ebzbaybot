package start

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/usvc/go-config"
	"github.com/zephinzer/ebzbaybot/internal/database"
	"github.com/zephinzer/ebzbaybot/internal/lifecycle"
	"github.com/zephinzer/ebzbaybot/internal/telegram"
	"github.com/zephinzer/ebzbaybot/internal/utils/log"
)

const (
	ConfigKeyMigrationsPath    = "migrations-path"
	ConfigKeyPostgresURL       = "postgres-url"
	ConfigKeyTelegramBotApiKey = "telegram-bot-api-key"
)

var conf = config.Map{
	ConfigKeyMigrationsPath: &config.String{
		Shorthand: "m",
		Default:   "data/migrations",
	},
	ConfigKeyPostgresURL: &config.String{
		Shorthand: "d",
		Default:   "postgres://user:password@localhost:35432/database?sslmode=disable",
	},
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
	// create & check database connection
	databaseUrl := conf.GetString(ConfigKeyPostgresURL)
	connection, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to the db: %s", err)
	}
	if err := connection.Ping(); err != nil {
		return fmt.Errorf("failed to ping the db: %s", err)
	}

	// migrate stuff
	migrationsPath := conf.GetString(ConfigKeyMigrationsPath)
	if err := database.Migrate(database.MigrateOpts{
		DatabaseURL:    databaseUrl,
		MigrationsPath: migrationsPath,
	}); err != nil {
		return fmt.Errorf("failed to migrate db: %s", err)
	}

	// start the bot
	botToken := conf.GetString(ConfigKeyTelegramBotApiKey)
	if botToken == "" {
		return fmt.Errorf("failed to receive a telegram bot token, pass it in with --%s", ConfigKeyTelegramBotApiKey)
	}

	var waiter sync.WaitGroup

	waiter.Add(1)

	go func() {
		lifecycleInterval := 10 * time.Second
		lifecycle.StartCollectionsScraping(lifecycle.ScrapingOpts{
			Connection:     connection,
			ScrapeInterval: lifecycleInterval,
		})
		waiter.Done()
	}()

	go func() {
		if err := telegram.StartBot(telegram.StartBotOpts{
			ApiKey:         botToken,
			Connection:     connection,
			IsDebugEnabled: false,
		}); err != nil {
			log.Errorf("failed to keep bot alive: %s", err)
			waiter.Done()
		}
	}()

	waiter.Wait()
	return nil
}
