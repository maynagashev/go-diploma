package app

import (
	"context"
	"fmt"
	"time"

	// Импортируем драйвер pgx для работы с PostgreSQL через database/sql.
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	maxOpenConns    = 25              // Максимальное количество открытых соединений
	connMaxLifetime = 5 * time.Minute // Максимальное время жизни соединения
	maxIdleConns    = 25              // Максимальное количество простаивающих соединений
	connMaxIdleTime = 5 * time.Minute // Максимальное время простоя соединения
)

// NewDB создает новое подключение к базе данных.
func NewDB(ctx context.Context, dsn string) (*sqlx.DB, error) {
	db, connectErr := sqlx.ConnectContext(ctx, "pgx", dsn)
	if connectErr != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", connectErr)
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime)

	// Проверка подключения
	if pingErr := db.PingContext(ctx); pingErr != nil {
		return nil, fmt.Errorf("не удалось проверить подключение к базе данных: %w", pingErr)
	}

	return db, nil
}
