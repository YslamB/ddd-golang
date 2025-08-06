package middleware

import (
	infra_logging "gddd/internal/infrastructure/logging"

	"github.com/gofiber/fiber/v2"
)

func Logger(zlog *infra_logging.Logger) fiber.Handler {

	return func(c *fiber.Ctx) error {
		err := c.Next()
		status := c.Response().StatusCode()

		if status != 200 {
			zlog.Infof("method: %s, path: %s, status: %s, Non-200 request error: %s", c.Method(), c.Path(), status, err.Error())
		}

		return err
	}
}
