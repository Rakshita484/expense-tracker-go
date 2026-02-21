package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/glebarez/sqlite"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds the application configuration.
type Config struct {
	DBDriver   string // "sqlite" or "postgres"
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	Port       string
}

// LoadConfig reads environment variables from .env and returns a Config struct.
func LoadConfig() *Config {
	// Load .env file if it exists; ignore error if it doesn't (e.g., in production).
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, reading configuration from environment")
	}

	return &Config{
		DBDriver:   getEnv("DB_DRIVER", "sqlite"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "expense_tracker"),
		Port:       getEnv("PORT", "8080"),
	}
}

// ConnectDatabase establishes a database connection using GORM.
// Supports both SQLite (zero-install, default) and PostgreSQL (production).
func ConnectDatabase(cfg *Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.DBDriver {
	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
		)
		dialector = postgres.Open(dsn)
		slog.Info("Using PostgreSQL database",
			"host", cfg.DBHost,
			"port", cfg.DBPort,
			"database", cfg.DBName,
		)

	default:
		// SQLite: zero configuration, stores data in a local file.
		dbFile := cfg.DBName + ".db"
		dialector = sqlite.Open(dbFile)
		slog.Info("Using SQLite database", "file", dbFile)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	slog.Info("Database connection established")
	return db, nil
}

// getEnv reads an environment variable or returns a default value.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
