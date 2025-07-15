package main

import (
	"ddd/configs"
	"ddd/internal/application/service"
	"ddd/internal/infrastructure/database"
	"ddd/internal/infrastructure/handler"
	"ddd/internal/infrastructure/repository"
	"ddd/internal/infrastructure/web"
	"ddd/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
)

func main() {

	appLogger := logger.New()

	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	db, err := database.NewPostgresDB(config.Database)
	if err != nil {
		appLogger.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)

	userService := service.NewUserService(userRepo)

	userHandler := handler.NewUserHandler(userService)

	app := fiber.New(fiber.Config{
		ErrorHandler: web.ErrorHandler,
	})

	web.SetupRoutes(app, userHandler)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		appLogger.Info("Shutting down server...")
		_ = app.Shutdown()
	}()

	appLogger.Infof("Server starting on port %s", config.Server.Port)
	log.Fatal(app.Listen(":" + config.Server.Port))
}
