package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gophermart/internal/handlers"
	"gophermart/internal/repository"
	"gophermart/internal/service"
	"gophermart/internal/worker"
)

const (
	defaultWorkerCount = 2
	defaultTimeout     = 10 * time.Second
)

// App представляет основную структуру приложения.
type App struct {
	echo           *echo.Echo
	db             *sqlx.DB
	userHandler    *handlers.UserHandler
	orderHandler   *handlers.OrderHandler
	balanceHandler *handlers.BalanceHandler
	accrualWorker  *worker.AccrualWorker
	config         Config
	wg             sync.WaitGroup // добавляем WaitGroup для ожидания завершения горутин
}

// New создает новый экземпляр приложения.
func New(ctx context.Context, cfg Config) (*App, error) {
	// Инициализация базы данных
	db, dbErr := NewDB(ctx, cfg.DatabaseURI)
	if dbErr != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", dbErr)
	}

	// Применяем миграции
	if migrateErr := MigrateDB(db.DB, cfg.MigrationsDir); migrateErr != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", migrateErr)
	}

	// Инициализация репозиториев
	userRepo := repository.NewUserRepo(db)
	orderRepo := repository.NewOrderRepo(db, slog.Default())
	balanceRepo := repository.NewBalanceRepo(db, slog.Default())

	// Инициализация сервисов
	userService := service.NewUserService(userRepo, cfg.JWTSecret, cfg.JWTExpirationPeriod)
	orderService := service.NewOrderService(orderRepo)
	balanceService := service.NewBalanceService(balanceRepo, slog.Default())
	accrualService := service.NewAccrualService(cfg.AccrualSystemAddress)

	// Создаем воркер для обработки начислений
	accrualWorker := worker.NewAccrualWorker(
		slog.Default(),
		orderRepo,
		accrualService,
		defaultWorkerCount, // количество воркеров
		defaultTimeout,
		0, // без задержки между попытками
	)

	// Инициализация обработчиков
	userHandler := handlers.NewUserHandler(userService)
	orderHandler := handlers.NewOrderHandler(orderService)
	balanceHandler := handlers.NewBalanceHandler(balanceService)

	// Инициализация Echo
	e := echo.New()
	e.Validator = NewValidator()

	// Промежуточное ПО (middleware)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	app := &App{
		echo:           e,
		db:             db,
		userHandler:    userHandler,
		orderHandler:   orderHandler,
		balanceHandler: balanceHandler,
		accrualWorker:  accrualWorker,
		config:         cfg,
	}

	// Настройка маршрутов
	app.setupRoutes()

	return app, nil
}

// Start запускает приложение.
func (a *App) Start(ctx context.Context, address string) error {
	// Запускаем воркер начислений в отдельной горутине
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		a.accrualWorker.Start(ctx)
	}()

	// Запускаем HTTP-сервер в фоне
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- a.echo.Start(address)
	}()

	// Ожидаем либо завершения контекста, либо ошибки сервера
	select {
	case <-ctx.Done():
		slog.Info("shutting down server...")

	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			return err
		}
	}

	// Создаём контекст с таймаутом для graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	if err := a.Shutdown(shutdownCtx); err != nil {
		slog.Error("failed to shutdown application", "error", err)
	}

	return nil
}

// Shutdown выполняет корректное завершение работы приложения.
func (a *App) Shutdown(ctx context.Context) error {
	// Ждем завершения всех воркеров
	shutdownComplete := make(chan struct{})
	go func() {
		a.wg.Wait()
		close(shutdownComplete)
	}()

	// Ожидаем либо завершения всех воркеров, либо таймаута контекста
	select {
	case <-shutdownComplete:
		slog.Info("all workers completed")
	case <-ctx.Done():
		slog.Warn("shutdown timeout exceeded, some workers may not have completed")
	}

	if err := a.db.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	if err := a.echo.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown http server: %w", err)
	}

	return nil
}

// setupRoutes настраивает маршруты приложения.
func (a *App) setupRoutes() {
	// Группа API
	api := a.echo.Group("/api")

	// Маршруты пользователя
	user := api.Group("/user")

	// Публичные маршруты
	user.POST("/register", a.userHandler.Register)
	user.POST("/login", a.userHandler.Authenticate)

	// Защищенные маршруты
	protected := user.Group("", JWTMiddleware(a.config.JWTSecret))

	// Маршруты заказов
	protected.POST("/orders", a.orderHandler.Register)
	protected.GET("/orders", a.orderHandler.GetOrders)

	// Маршруты баланса
	protected.GET("/balance", a.balanceHandler.GetBalance)
	protected.POST("/balance/withdraw", a.balanceHandler.Withdraw)
	protected.GET("/withdrawals", a.balanceHandler.GetWithdrawals)
}
