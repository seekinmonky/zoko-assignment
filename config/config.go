package config

import (
	"fmt"
	"os"
)

// Config holds the application configuration
type Config struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	RedisAddr        string
	RedisPassword    string
	ServerPort       string
}

// GetPostgresConnStr returns the PostgreSQL connection string
func (c *Config) GetPostgresConnStr() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPassword, c.PostgresDB)
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:     getEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
		PostgresDB:       getEnv("POSTGRES_DB", "messages"),
		RedisAddr:        getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:    getEnv("REDIS_PASSWORD", ""),
		ServerPort:       getEnv("SERVER_PORT", "8080"),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
