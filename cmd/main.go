package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"avito-shop/config"
	"avito-shop/db"
	"avito-shop/internal/app"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	if err := cfg.InitDB(); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}

	if err := db.Migrate(cfg.DB); err != nil {
		log.Fatalf("Ошибка выполнения миграций базы данных: %v", err)
	}

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Ошибка создания приложения: %v", err)
	}

	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: application.Router(),
	}

	// Запуск сервера в горутине
	go func() {
		log.Printf("Запуск сервера на адресе %s", cfg.ServerAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	// Ожидание сигнала для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Начало завершения работы сервера...")

	// Контекст для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown сервера
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Принудительное завершение работы сервера: %v", err)
	}

	// Очистка ресурсов приложения
	if err := application.Cleanup(); err != nil {
		log.Fatalf("Ошибка очистки ресурсов приложения: %v", err)
	}

	log.Println("Сервер успешно завершил работу")
}
