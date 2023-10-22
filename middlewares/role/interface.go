package role

import "github.com/gofiber/fiber/v2"

type Middleware interface {
	Authorization(allowedRoles ...string) fiber.Handler
}
