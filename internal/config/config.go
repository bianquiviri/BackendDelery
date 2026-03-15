package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBSSLMode  string
	Port       string
	GinMode    string
}

// LoadConfig reads the .env file and populates the Config struct.
// It applies defense in depth by providing defaults or failing fast if critical vars are missing.
func LoadConfig() *Config {
	// Ignore error if .env doesn't exist, we might be using actual system env vars (e.g. in Docker/CI)
	_ = godotenv.Load()

	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "secret"),
		DBName:     getEnv("DB_NAME", "daas_db"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		Port:       getEnv("PORT", "8084"),
		GinMode:    getEnv("GIN_MODE", "release"),
	}

	// Validate critical config
	if cfg.DBPassword == "" {
		log.Println("WARNING: DB_PASSWORD is not set. This might cause connection failures.")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
