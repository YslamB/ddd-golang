package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		// It's okay if .env file doesn't exist, load as default values
	}

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT"),
			Env:  getEnv("ENV"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST"),
			Port:     getEnv("DB_PORT"),
			User:     getEnv("DB_USER"),
			Password: getEnv("DB_PASSWORD"),
			DBName:   getEnv("DB_NAME"),
			SSLMode:  getEnv("DB_SSLMODE"),
		},
	}

	return config, nil
}

func getEnv(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Fatalf("%s: Not in .env file", key)
	}

	return value
}
