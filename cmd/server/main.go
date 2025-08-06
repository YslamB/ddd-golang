package main

import (
	"context"
	"gddd/internal/infrastructure/config"
	infra_logging "gddd/internal/infrastructure/logging"
	infra_store "gddd/internal/infrastructure/storage"
	infra_web "gddd/internal/infrastructure/web"
	"gddd/internal/shared/utils"
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
	ctxTimeOut, cancel := context.WithTimeout(context.Background(), cfg.Listen.DBCtxTimeout)
	psqlConn := infra_store.PostgresInit(&ctxTimeOut, cfg, log)
	defer cancel()
	defer psqlConn.Close()
	fiberConfig := utils.FiberConfig(cfg)
	app := fiber.New(fiberConfig)
	infra_web.SetupRoutes(app, psqlConn, log)

	go func() {
		log.Infof("Server starting on port: %s", cfg.Listen.Port)
		err := app.Listen(":" + cfg.Listen.Port)

		if err != nil && err != http.ErrServerClosed {
			log.Errorf("error server close: %s", err.Error())
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
