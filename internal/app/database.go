package app

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	maxOpenConns    = 25              // Максимальное количество открытых соединений
	connMaxLifetime = 5 * time.Minute // Максимальное время жизни соединения
	maxIdleConns    = 25              // Максимальное количество простаивающих соединений
	connMaxIdleTime = 5 * time.Minute // Максимальное время простоя соединения
)

// NewDB создает новое подключение к базе данных
func NewDB(ctx context.Context, dsn string) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, "pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime)

	// Проверка подключения
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("не удалось проверить подключение к базе данных: %w", err)
	}

	return db, nil
}
