package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	loggerFiber "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/edwynrrangel/tasks/config"
	"github.com/edwynrrangel/tasks/logger"
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

	logger.Info("Starting application")
	logger.Fatal(
		app.Listen(
			fmt.Sprintf(":%s", config.Config.Port),
		),
	)
}
