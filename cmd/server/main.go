package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	app_product "gddd/internal/application/product"
	app_user "gddd/internal/application/user"
	domain_user "gddd/internal/domain/user"
	infra_config "gddd/internal/infrastructure/config"
	infra_logging "gddd/internal/infrastructure/logging"
	infra_inmemory "gddd/internal/infrastructure/persistence/inmemory" // For in-memory repos
	infra_postgres "gddd/internal/infrastructure/persistence/postgres" // For Postgres repos
	http_interfaces "gddd/internal/interfaces/http"

	fiber_recover "github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Initialize custom logger
	appLogger := infra_logging.NewLogger(os.Stdout)

	// Load configuration
	cfg := infra_config.LoadConfig()

	// Initialize Database (PostgreSQL example)
	// In a real app, handle connection pooling, retries, etc.
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		appLogger.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			appLogger.Errorf("Error closing database connection: %v", err)
		}
	}()

	if err = db.Ping(); err != nil {
		appLogger.Fatalf("Failed to ping database: %v", err)
	}
	appLogger.Infof("Successfully connected to database!")

	// Choose repository implementation based on config or environment
	var userRepo domain_user.Repository
	var productRepo app_product.ProductRepository

	if cfg.UseInMemoryRepos {
		appLogger.Infof("Using in-memory repositories.")
		userRepo = infra_inmemory.NewInMemoryUserRepository()
		productRepo = infra_inmemory.NewInMemoryProductRepository()
	} else {
		appLogger.Infof("Using PostgreSQL repositories.")
		userRepo = infra_postgres.NewPostgresUserRepository(db)
		productRepo = infra_postgres.NewPostgresProductRepository(db)
	}

	// Initialize Application Services
	userService := app_user.NewService(userRepo, appLogger)
	productService := app_product.NewService(productRepo, appLogger)
	customErrorHandler := func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		// Log the error using your custom logger
		appLogger.Errorf("Fiber Error: %v, Path: %s, IP: %s", err, c.Path(), c.IP())

		// Send custom error response
		return c.Status(code).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// Initialize Fiber App
	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.ServerReadTimeout,
		WriteTimeout: cfg.ServerWriteTimeout,
		IdleTimeout:  cfg.ServerIdleTimeout,
		ErrorHandler: customErrorHandler,
	})

	// Add Fiber Middleware
	app.Use(fiber_recover.New())      // Recovers from panics and sends 500
	app.Use(logger.New(logger.Config{ // Request logger
		Format: "[${time}] ${ip} ${status} - ${method} ${path} ${latency}\n",
		Output: os.Stdout,
	}))

	// Initialize HTTP Handlers
	userHandler := http_interfaces.NewUserHandler(userService, appLogger)
	productHandler := http_interfaces.NewProductHandler(productService, appLogger)

	// Setup HTTP Routes
	apiV1 := app.Group("/api/v1")
	{
		// User routes
		apiV1.Post("/auth/register", userHandler.RegisterUser)
		apiV1.Get("/users/:id", userHandler.GetUserByID)
		apiV1.Put("/users/:id", userHandler.UpdateUser)
		apiV1.Delete("/users/:id", userHandler.DeleteUser)

		// Product routes
		apiV1.Post("/products", productHandler.CreateProduct)
		apiV1.Get("/products/:id", productHandler.GetProductByID)
		apiV1.Put("/products/:id", productHandler.UpdateProduct)
		apiV1.Delete("/products/:id", productHandler.DeleteProduct)
	}

	// Start server in a goroutine
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = cfg.ServerPort // Use configured port
		}
		appLogger.Infof("Server starting on port %s", port)
		if err := app.Listen(":" + port); err != nil {
			appLogger.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful Shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // Listen for OS signals
	<-c                                             // Block until a signal is received
	appLogger.Infof("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Give 10 seconds for graceful shutdown
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		appLogger.Fatalf("Server forced to shutdown: %v", err)
	}
	appLogger.Infof("Server exited gracefully.")
}
