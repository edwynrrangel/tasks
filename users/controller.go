package users

import (
	"net/http"

	"github.com/edwynrrangel/tasks/errors"
	"github.com/edwynrrangel/tasks/logger"
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

func (c *controller) Create(ctx *fiber.Ctx) error {

	var body = new(CreateRequest)
	if err := validator.BodyParser(ctx, body); err != nil {
		logger.Error(
			"error parsing body",
			"func", "Login - ctx.BodyParser",
			"error", err.Error(),
		)
		return errors.New(http.StatusBadRequest, err.Error(), nil).
			ToFiber(ctx)
	}

	response, err := c.service.Create(body)
	if err != nil {
		return err.ToFiber(ctx)
	}

	return ctx.Status(http.StatusCreated).JSON(response)
}

func (c *controller) List(ctx *fiber.Ctx) error {

	var queryParams = new(ListRequest)
	if err := validator.QueryParser(ctx, queryParams); err != nil {
		logger.Error(
			"error parsing query params",
			"func", "Login - ctx.QueryParser",
			"error", err.Error(),
		)
		return errors.New(http.StatusBadRequest, err.Error(), nil).
			ToFiber(ctx)
	}

	response, err := c.service.List(*queryParams)
	if err != nil {
		return err.ToFiber(ctx)
	}

	return ctx.Status(http.StatusOK).JSON(response)
}
