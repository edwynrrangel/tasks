package auth

import (
	"github.com/gofiber/fiber/v2"

	"github.com/edwynrrangel/tasks/errors"
	"github.com/edwynrrangel/tasks/middlewares/session"
)

type Service interface {
	Login(body LoginRequest) (*LoginResponse, *errors.Error)
	ChangePassword(authUser *session.Auth, body ChangePasswordRequest) (*LoginResponse, *errors.Error)
	Logout(authUser *session.Auth) *errors.Error
}

type Controller interface {
	Login(ctx *fiber.Ctx) error
	ChangePassword(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
}
