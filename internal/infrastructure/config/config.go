package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration settings.
type Config struct {
	ServerPort         string
	ServerReadTimeout  time.Duration
	ServerWriteTimeout time.Duration
	ServerIdleTimeout  time.Duration
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	UseInMemoryRepos   bool // For local development/testing
}

// LoadConfig loads configuration from environment variables.
// In a real application, you might use a library like Viper for more robust config management.
func LoadConfig() *Config {
	port := getEnv("PORT", "8001")
	readTimeoutStr := getEnv("SERVER_READ_TIMEOUT", "5s")
	writeTimeoutStr := getEnv("SERVER_WRITE_TIMEOUT", "10s")
	idleTimeoutStr := getEnv("SERVER_IDLE_TIMEOUT", "120s")
	useInMemoryStr := getEnv("USE_IN_MEMORY_REPOS", "false")

	readTimeout, err := time.ParseDuration(readTimeoutStr)
	if err != nil {
		log.Fatalf("Invalid SERVER_READ_TIMEOUT: %v", err)
	}
	writeTimeout, err := time.ParseDuration(writeTimeoutStr)
	if err != nil {
		log.Fatalf("Invalid SERVER_WRITE_TIMEOUT: %v", err)
	}
	idleTimeout, err := time.ParseDuration(idleTimeoutStr)
	if err != nil {
		log.Fatalf("Invalid SERVER_IDLE_TIMEOUT: %v", err)
	}

	useInMemory, err := strconv.ParseBool(useInMemoryStr)
	if err != nil {
		log.Fatalf("Invalid USE_IN_MEMORY_REPOS: %v", err)
	}

	return &Config{
		ServerPort:         port,
		ServerReadTimeout:  readTimeout,
		ServerWriteTimeout: writeTimeout,
		ServerIdleTimeout:  idleTimeout,
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", "password"),
		DBName:             getEnv("DB_NAME", "your_database"),
		UseInMemoryRepos:   useInMemory,
	}
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
