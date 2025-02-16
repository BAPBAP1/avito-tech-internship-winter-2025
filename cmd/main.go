package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	//"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"

	"github.com/BAPBAP1/avito-tech-internship-winter-2025/config"
	"github.com/BAPBAP1/avito-tech-internship-winter-2025/internal/handler"
	"github.com/BAPBAP1/avito-tech-internship-winter-2025/internal/middleware"
	"github.com/BAPBAP1/avito-tech-internship-winter-2025/internal/repository"
	"github.com/BAPBAP1/avito-tech-internship-winter-2025/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := runMigrations(db, cfg); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	merchRepo := repository.NewMerchRepository()

	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	walletService := service.NewWalletService(userRepo, transactionRepo, db)
	merchService := service.NewMerchService(merchRepo, userRepo, db)

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	authHandler := handler.NewAuthHandler(authService)
	r.POST("/auth", authHandler.Login)

	authMiddleware := middleware.JWTAuthMiddleware(cfg.JWTSecret)

	authorized := r.Group("/api")
	authorized.Use(authMiddleware)
	{
		walletHandler := handler.NewWalletHandler(walletService)
		authorized.POST("/transfer", walletHandler.Transfer)
		authorized.GET("/wallet", walletHandler.GetWallet)
		authorized.GET("/wallet/history", walletHandler.GetWalletHistory)

		merchHandler := handler.NewMerchHandler(merchService)
		authorized.GET("/merch", merchHandler.ListMerch)
		authorized.POST("/purchase", merchHandler.PurchaseMerch)
		authorized.GET("/purchases", merchHandler.ListPurchases)
	}

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Server starting on port 8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func runMigrations(db *sql.DB, cfg *config.Config) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		cfg.DBName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if os.Getenv("MIGRATE_DROP") == "true" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("failed to drop database: %w", err)
		}
		log.Println("Migrations dropped")
		return nil
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations up: %w", err)
	}
	log.Println("Migrations successfully applied")
	return nil
}
