package app

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

// MigrateDB применяет миграции базы данных.
func MigrateDB(db *sql.DB, migrationsDir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("не удалось установить диалект базы данных: %w", err)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("не удалось применить миграции: %w", err)
	}

	return nil
}
