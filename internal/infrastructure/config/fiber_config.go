package config

import (
	"gddd/internal/shared/utils"

	"github.com/gofiber/fiber/v2"
)

func FiberConfig(cfg *Config) fiber.Config {
	return fiber.Config{
		WriteTimeout: cfg.Listen.WriteTimeout,
		ReadTimeout:  cfg.Listen.ReadTimeout,
		IdleTimeout:  cfg.Listen.IDLETimeout,
		ErrorHandler: utils.ErrorHandler,
	}
}
