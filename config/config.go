package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"avito-shop/db"
)

type Config struct {
	DB            *gorm.DB
	DatabaseHost  string
	DatabaseUser  string
	DatabasePass  string
	DatabaseName  string
	DatabasePort  string
	ServerAddress string
	JWTSecret     string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("ошибка загрузки файла .env: %w", err)
	}

	config := &Config{
		DatabaseHost:  getEnv("DATABASE_HOST", "localhost"),
		DatabaseUser:  getEnv("DATABASE_USER", "postgres"),
		DatabasePass:  getEnv("DATABASE_PASSWORD", "password"),
		DatabaseName:  getEnv("DATABASE_NAME", "shop"),
		DatabasePort:  getEnv("DATABASE_PORT", "5432"),
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
	}

	return config, nil
}

func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.DatabaseHost,
		c.DatabaseUser,
		c.DatabasePass,
		c.DatabaseName,
		c.DatabasePort,
	)
}

func (c *Config) InitDB() error {
	dsn := c.GetDatabaseDSN()
	database, err := db.NewDB(dsn)
	if err != nil {
		return fmt.Errorf("не удалось инициализировать базу данных: %w", err)
	}
	c.DB = database
	return nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
