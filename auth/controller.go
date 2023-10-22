package auth

import (
	"net/http"

	"github.com/edwynrrangel/tasks/errors"
	"github.com/edwynrrangel/tasks/logger"
	"github.com/edwynrrangel/tasks/middlewares/session"
	"github.com/edwynrrangel/tasks/validator"
	"github.com/gofiber/fiber/v2"
)

type controller struct {
	service Service
}

func NewController(service Service) Controller {
	return &controller{
		service: service,
	}
}

func (c *controller) Login(ctx *fiber.Ctx) error {

	var body = new(LoginRequest)
	if err := validator.BodyParser(ctx, body); err != nil {
		logger.Error(
			"error parsing body",
			"func", "Login - ctx.BodyParser",
			"error", err.Error(),
		)
		return errors.New(http.StatusBadRequest, err.Error(), nil).
			ToFiber(ctx)
	}

	response, err := c.service.Login(*body)
	if err != nil {
		logger.Error(
			"error login",
			"func", "Login - service.Login",
			"error", err.Message,
		)
		return err.ToFiber(ctx)
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

func (c *controller) ChangePassword(ctx *fiber.Ctx) error {

	var body = new(ChangePasswordRequest)
	if err := validator.BodyParser(ctx, body); err != nil {
		logger.Error(
			"error parsing body",
			"func", "ChangePassword - ctx.BodyParser",
			"error", err.Error(),
		)
		return errors.New(http.StatusBadRequest, err.Error(), nil).
			ToFiber(ctx)
	}

	authUser := ctx.Locals("user").(*session.Auth)
	response, err := c.service.ChangePassword(authUser, *body)
	if err != nil {
		logger.Error(
			"error change password",
			"func", "ChangePassword - service.ChangePassword",
			"error", err.Message,
		)
		return err.ToFiber(ctx)
	}
	return ctx.Status(http.StatusOK).JSON(response)
}

func (c *controller) Logout(ctx *fiber.Ctx) error {
	authUser := ctx.Locals("user").(*session.Auth)
	err := c.service.Logout(authUser)
	if err != nil {
		logger.Error(
			"error logout",
			"func", "Logout - service.Logout",
			"error", err.Message,
		)
		return err.ToFiber(ctx)
	}
	return ctx.SendStatus(http.StatusNoContent)
}
