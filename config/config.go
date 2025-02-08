package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	JwtSecretKey string
	AppPort      string
}

func LoadConfig() *Config {
	// Загружаем переменные из .env файла
	err := godotenv.Load()
	if err != nil {
		log.Printf("Using environment variables from docker or os")

	} else {
		log.Printf("Using .env file")
	}

	// Возвращаем структуру с конфигурацией
	return &Config{
		DBHost:       getEnv("DB_HOST", "auction-db"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPassword:   getEnv("DB_PASSWORD", "password"),
		DBName:       getEnv("DB_NAME", "auction_db"),
		JwtSecretKey: getEnv("JWT_SECRET_KEY", "your_default_jwt_secret_key"),
		AppPort:      getEnv("APP_PORT", "8000"),
	}
}

func (c *Config) GetDBConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
