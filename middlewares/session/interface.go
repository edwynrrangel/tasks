package session

import "github.com/gofiber/fiber/v2"

type Middleware interface {
	Validate(ctx *fiber.Ctx) error
	ValidateFirstLogin(ctx *fiber.Ctx) error
}
