package app

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"avito-shop/config"
	"avito-shop/internal/handlers"
	"avito-shop/internal/middleware"
	"avito-shop/internal/ports"
	"avito-shop/internal/repositories"
	"avito-shop/internal/services"
)

type App struct {
	router           *gin.Engine
	userService      ports.UserService
	authService      ports.AuthService
	inventoryService ports.InventoryService
	userHandler      ports.UserHandler
	authHandler      ports.AuthHandler
	inventoryHandler ports.InventoryHandler
	db               *sql.DB
}

func New(cfg *config.Config) (*App, error) {
	// Создаем репозитории
	userRepo := repositories.NewUserRepository(cfg.DB)
	inventoryRepo := repositories.NewInventoryRepository(cfg.DB)
	itemRepo := repositories.NewItemRepository(cfg.DB)
	transactionRepo := repositories.NewTransactionRepository(cfg.DB)

	// Создаем сервисы
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(cfg.DB, userRepo, inventoryRepo, transactionRepo)
	inventoryService := services.NewInventoryService(cfg.DB, inventoryRepo, userRepo, itemRepo)

	// Создаем обработчики
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)

	// Создаем роутер
	router := gin.Default()

	// Настройка middleware
	authMiddleware := middleware.NewAuthMiddleware()

	// Публичные роуты
	router.POST("/api/auth", authHandler.Login)

	// Защищенные роуты
	authorized := router.Group("/")
	authorized.Use(authMiddleware.Handle())
	{
		authorized.GET("/api/info", userHandler.GetUserInfo)
		authorized.POST("/api/sendCoin", userHandler.SendCoins)
		authorized.GET("/api/buy/:item", inventoryHandler.BuyItem)
	}

	sqlDB, err := cfg.DB.DB()
	if err != nil {
		log.Printf("ошибка при получении соединения с БД: %v", err)
		return nil, err
	}

	if sqlDB == nil {
		log.Println("ошибка: sqlDB is nil")
		return nil, fmt.Errorf("ошибка: sqlDB is nil")
	}

	return &App{
		router:           router,
		userService:      userService,
		authService:      authService,
		inventoryService: inventoryService,
		userHandler:      userHandler,
		authHandler:      authHandler,
		inventoryHandler: inventoryHandler,
		db:               sqlDB,
	}, nil
}

func (a *App) Router() *gin.Engine {
	return a.router
}

func (a *App) Cleanup() error {
	if a.db != nil {
		err := a.db.Close()
		if err != nil {
			log.Printf("ошибка при закрытии соединения с БД: %v", err)
			return err
		}
		log.Println("Соединение с БД успешно закрыто")
		return nil
	}
	return nil
}
