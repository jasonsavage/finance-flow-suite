package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    string
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
}

func Load() *Config {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	return &Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "financeflow"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
