package main

import (
	"gddd/internal/infrastructure/config"
	infra_logging "gddd/internal/infrastructure/logging"
	"gddd/internal/infrastructure/storage"
	"gddd/internal/infrastructure/web"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.Init()

	log := infra_logging.InitLogger(cfg)
	psqlConn := storage.PostgresInit(cfg)

	defer psqlConn.Close()

	app := fiber.New(fiber.Config{
		WriteTimeout: cfg.Listen.WriteTimeout,
		ReadTimeout:  cfg.Listen.ReadTimeout,
		IdleTimeout:  cfg.Listen.IDLETimeout,
		ErrorHandler: web.ErrorHandler,
	})

	web.SetupRoutes(app, psqlConn, log)

	go func() {
		log.Infof("Server starting on port: %s", cfg.Listen.Port)
		err := app.Listen(":" + cfg.Listen.Port)

		if err != nil && err != http.ErrServerClosed {
			log.Errorf(err.Error())
		}
	}()

	// When a shutdown signal (like Ctrl+C) is caught, you initiate a graceful shutdown using srv.Shutdown(ctx), giving active connections time to complete before the server shuts down.
	// Graceful shutdown is handled within the main function, allowing ongoing requests to finish processing before the server shuts down, if req not completed in 5 seconds the server will force shutdown.
	// If the expected request finishes before 5 seconds, the server is shut down immediately.
	// New Requests will not be accepted.
	// Wait for an interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infof("Shutting down server...")

	shutdownCh := make(chan struct{})
	go func() {
		if err := app.Shutdown(); err != nil {
			log.Errorf("Server forced to shutdown:%s", err.Error())
		}
		close(shutdownCh)
	}()

	select {
	case <-shutdownCh:
		log.Infof("Server shutdown gracefully")
	case <-time.After(5 * time.Second):
		log.Errorf("Server shutdown timed out, forcing exit")
	}

}
