package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	DBType   string
	PGURL    string
	RedisURL string
}

func Load() (*Config, error) {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file. default values will be used")
	}

	cfg := &Config{
		Port:     getEnv("APP_PORT", "8080"),
		DBType:   getEnv("DB_TYPE", "inmemory"),
		PGURL:    getEnv("PG_URL", ""),
		RedisURL: getEnv("REDIS_URL", "redis://localhost:6379"),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)

	if !exists {
		return defaultValue
	}
	return value
}
