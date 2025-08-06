package web

import (
	app_user "gddd/internal/application/user"
	infra_logging "gddd/internal/infrastructure/logging"
	infra_postgres "gddd/internal/infrastructure/persistence/postgres" // For Postgres repos
	"gddd/internal/infrastructure/web/middleware"
	http_interfaces "gddd/internal/interfaces/http"
	db "gddd/pkg/database"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, db *db.Database, log *infra_logging.Logger) {
	userRepo := infra_postgres.NewPostgresUserRepository(db)
	userService := app_user.NewService(userRepo, log)
	userHandler := http_interfaces.NewUserHandler(userService, log)

	app.Use(middleware.Logger(log))
	app.Use(middleware.CORS())

	api := app.Group("/api/v1")
	usersAPI := api.Group("/users")
	SetupUserRoutes(usersAPI, userHandler)
}

func SetupUserRoutes(r fiber.Router, handler *http_interfaces.UserHandler) {
	r.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	r.Put("/:id", handler.UpdateUser)
	r.Delete("/:id", handler.DeleteUser)
}
