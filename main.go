package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	loggerFiber "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/edwynrrangel/tasks/auth"
	"github.com/edwynrrangel/tasks/config"
	"github.com/edwynrrangel/tasks/logger"
	"github.com/edwynrrangel/tasks/tasks"
	"github.com/edwynrrangel/tasks/users"
)

func init() {
	config.Config = config.LoadConfig()
}

func main() {
	defer logger.Sync()

	app := fiber.New()
	app.Use(loggerFiber.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: config.Config.CorsAllowedOrigins,
	}))

	auth.ApplyRoutes(app, config.Config)
	users.ApplyRoutes(app, config.Config)
	tasks.ApplyRoutes(app, config.Config)
	app.Use(
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "lo sentimos, no se encontró el recurso solicitado",
			})
		},
	)

	logger.Info("Starting application")
	logger.Fatal(
		app.Listen(
			fmt.Sprintf(":%s", config.Config.Port),
		),
	)
}
