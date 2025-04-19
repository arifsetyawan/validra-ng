package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Type     string // Type will always be "postgres" now
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// Load loads configuration from environment variables
// It first attempts to load from a .env file if it exists
func Load() *Config {
	// Try to load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// It's okay if the .env file does not exist, we'll use environment variables
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Server: ServerConfig{
			Port:         getEnvAsInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvAsInt("SERVER_READ_TIMEOUT", 60),
			WriteTimeout: getEnvAsInt("SERVER_WRITE_TIMEOUT", 60),
		},
		Database: DatabaseConfig{
			Type:     getEnv("DB_TYPE", "postgres"),
			Host:     getEnv("DB_HOST", "postgres"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "validra"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
	}
}

// getEnv retrieves environment variables or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves environment variables as int or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
