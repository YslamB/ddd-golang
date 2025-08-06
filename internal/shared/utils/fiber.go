package utils

import (
	infra_config "gddd/internal/infrastructure/config"

	"github.com/gofiber/fiber/v2"
)

func FiberConfig(cfg *infra_config.Config) fiber.Config {
	return fiber.Config{
		WriteTimeout: cfg.Listen.WriteTimeout,
		ReadTimeout:  cfg.Listen.ReadTimeout,
		IdleTimeout:  cfg.Listen.IDLETimeout,
		ErrorHandler: ErrorHandler,
	}
}
