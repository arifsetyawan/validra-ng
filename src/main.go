package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arifsetyawan/validra/src/config"
	"github.com/arifsetyawan/validra/src/internal/delivery/http/handler"
	"github.com/arifsetyawan/validra/src/internal/delivery/http/middleware"
	"github.com/arifsetyawan/validra/src/internal/domain"
	"github.com/arifsetyawan/validra/src/internal/repository"
	"github.com/arifsetyawan/validra/src/internal/service"
	"github.com/arifsetyawan/validra/src/pkg/database"
	"github.com/arifsetyawan/validra/src/pkg/logger"
	"github.com/arifsetyawan/validra/src/pkg/validator"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	// Import generated docs
	_ "github.com/arifsetyawan/validra/docs"
)

// @title Validra Engine API
// @version 1.0
// @description API documentation for the Validra Engine service
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.validra.io/support
// @contact.email support@validra.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	// Initialize logger
	log := logger.NewLogger()
	log.Info("Starting Validra Engine application...")

	// Load configuration
	cfg := config.Load()
	log.Info("Configuration loaded")
	// Initialize database based on type
	var resourceRepo domain.ResourceRepository

	if cfg.Database.Type == "postgres" {
		// Initialize PostgreSQL with GORM
		db, err := database.NewPostgresDB(
			cfg.Database.Host,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Name,
			cfg.Database.Port,
			cfg.Database.SSLMode,
		)
		if err != nil {
			log.Error("Failed to connect to PostgreSQL database: %v", err)
			os.Exit(1)
		}
		defer db.Close()

		// Run migrations
		if err := db.Migrate(); err != nil {
			log.Error("Failed to migrate database: %v", err)
			os.Exit(1)
		}
		log.Info("PostgreSQL database migrations completed")

		// Initialize repository with GORM
		resourceRepo = repository.NewGormResourceRepository(db)

	} else {
		// Initialize SQLite (for backwards compatibility)
		db, err := database.NewSQLiteDB(cfg.Database.Path)
		if err != nil {
			log.Error("Failed to connect to SQLite database: %v", err)
			os.Exit(1)
		}
		defer db.Close()

		// Run migrations
		if err := db.Migrate(); err != nil {
			log.Error("Failed to migrate SQLite database: %v", err)
			os.Exit(1)
		}
		log.Info("SQLite database migrations completed")

		// Initialize repository with SQLite
		resourceRepo = repository.NewSQLiteResourceRepository(db)
	}

	// Initialize Echo
	e := echo.New()
	e.Validator = validator.NewCustomValidator()

	// Setup middleware
	middleware.SetupMiddleware(e, log)

	// Initialize services
	resourceService := service.NewResourceService(resourceRepo)

	// Initialize handlers
	resourceHandler := handler.NewResourceHandler(resourceService)

	// Register routes
	resourceHandler.Register(e)

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Setup Swagger
	e.GET("/docs", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/docs/index.html")
	})
	e.GET("/docs/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/docs/index.html")
	})
	e.GET("/docs/*", echoSwagger.WrapHandler)
	log.Info("Swagger documentation available at /docs")

	// Start server
	go func() {
		address := fmt.Sprintf(":%d", cfg.Server.Port)
		log.Info("Server starting on port %d", cfg.Server.Port)

		s := &http.Server{
			Addr:         address,
			ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		}

		if err := e.StartServer(s); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server: %v", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown: %v", err)
		os.Exit(1)
	}

	log.Info("Server stopped gracefully")
}
