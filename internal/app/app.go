package app

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gophermart/internal/handlers"
	"gophermart/internal/repository"
	"gophermart/internal/service"
)

// App представляет основную структуру приложения
type App struct {
	echo        *echo.Echo
	db          *sqlx.DB
	userHandler *handlers.UserHandler
	config      Config
}

// New создает новый экземпляр приложения
func New(ctx context.Context, cfg Config) (*App, error) {
	// Инициализация базы данных
	db, err := NewDB(ctx, cfg.DatabaseURI)
	if err != nil {
		return nil, err
	}

	// Применение миграций
	if err := MigrateDB(db.DB, cfg.MigrationsDir); err != nil {
		return nil, err
	}

	// Инициализация репозиториев
	userRepo := repository.NewUserRepo(db)

	// Инициализация сервисов
	userService := service.NewUserService(userRepo, cfg.JWTSecret, cfg.JWTExpirationPeriod)

	// Инициализация обработчиков
	userHandler := handlers.NewUserHandler(userService)

	// Инициализация Echo
	e := echo.New()
	e.Validator = NewValidator()

	// Промежуточное ПО (middleware)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	app := &App{
		echo:        e,
		db:          db,
		userHandler: userHandler,
		config:      cfg,
	}

	// Настройка маршрутов
	app.setupRoutes()

	return app, nil
}

// Start запускает приложение
func (a *App) Start(address string) error {
	return a.echo.Start(address)
}

// Shutdown выполняет корректное завершение работы приложения
func (a *App) Shutdown(ctx context.Context) error {
	if err := a.db.Close(); err != nil {
		return err
	}
	return a.echo.Shutdown(ctx)
}

// setupRoutes настраивает маршруты приложения
func (a *App) setupRoutes() {
	// Группа API
	api := a.echo.Group("/api")

	// Маршруты пользователя
	user := api.Group("/user")

	// Публичные маршруты
	user.POST("/register", a.userHandler.Register)
	user.POST("/login", a.userHandler.Authenticate)

	// Защищенные маршруты будут добавлены позже
	protected := user.Group("", JWTMiddleware(a.config.JWTSecret))
	_ = protected // временно, чтобы избежать ошибки неиспользуемой переменной
}
