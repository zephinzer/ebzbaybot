package database

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type MigrateOpts struct {
	MigrationsPath string
	DatabaseURL    string
}

func Migrate(opts MigrateOpts) error {
	migrations, err := migrate.New(
		fmt.Sprintf("file://%s", opts.MigrationsPath),
		opts.DatabaseURL,
	)
	if err != nil {
		return fmt.Errorf("failed to init migrations: %s", err)
	}
	if err := migrations.Up(); err != nil {
		return fmt.Errorf("failed to run migrations: %s", err)
	}
	return nil
}
