package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/raksh/expense-tracker/config"
	"github.com/raksh/expense-tracker/handlers"
	"github.com/raksh/expense-tracker/models"
	"github.com/raksh/expense-tracker/repository"
	"github.com/raksh/expense-tracker/routes"
	"github.com/raksh/expense-tracker/services"
)

func main() {
	// Initialize structured logger for the entire application.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting Expense Tracker API")

	// Load configuration from .env file and environment variables.
	cfg := config.LoadConfig()

	// Connect to database (SQLite or PostgreSQL) via GORM.
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err.Error())
		os.Exit(1)
	}

	// Auto-migrate database schema.
	slog.Info("Running database migrations")
	if err := db.AutoMigrate(
		&models.User{},
		&models.Group{},
		&models.Expense{},
		&models.ExpenseSplit{},
	); err != nil {
		slog.Error("Failed to run migrations", "error", err.Error())
		os.Exit(1)
	}
	slog.Info("Database migrations completed")

	// Wire up dependencies.
	repo := repository.NewRepository(db)
	svc := services.NewService(repo)
	handler := handlers.NewHandler(svc)

	// Initialize Gin router.
	router := gin.Default()

	// Configure CORS for the frontend.
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Register all API routes.
	routes.RegisterRoutes(router, handler)

	// Start the HTTP server.
	addr := fmt.Sprintf(":%s", cfg.Port)
	slog.Info("Server starting", "address", addr)
	if err := router.Run(addr); err != nil {
		slog.Error("Server failed to start", "error", err.Error())
		os.Exit(1)
	}
}
