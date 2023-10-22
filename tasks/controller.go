package tasks

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/edwynrrangel/tasks/errors"
	"github.com/edwynrrangel/tasks/logger"
	"github.com/edwynrrangel/tasks/middlewares/session"
	"github.com/edwynrrangel/tasks/validator"
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
	var query = new(ListRequest)
	if err := validator.QueryParser(ctx, query); err != nil {
		logger.Error(
			"error parsing query",
			"func", "List - ctx.QueryParser",
			"error", err.Error(),
		)
		return errors.New(http.StatusBadRequest, err.Error(), nil).
			ToFiber(ctx)
	}

	response, err := c.service.List(ctx.Locals("user").(*session.Auth), *query)
	if err != nil {
		return err.ToFiber(ctx)
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

func (c *controller) Update(ctx *fiber.Ctx) error {
	var body = new(UpdateRequest)
	if err := validator.BodyParser(ctx, body); err != nil {
		logger.Error(
			"error parsing body",
			"func", "Update - ctx.BodyParser",
			"error", err.Error(),
		)
		return errors.New(http.StatusBadRequest, err.Error(), nil).
			ToFiber(ctx)
	}

	// Update status
	if body.Status != "" {
		response, err := c.service.UpdateStatus(ctx.Params("id"), body)
		if err != nil {
			return err.ToFiber(ctx)
		}

		return ctx.Status(http.StatusOK).JSON(response)
	}

	// Update task
	response, err := c.service.Update(ctx.Params("id"), body)
	if err != nil {
		return err.ToFiber(ctx)
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

func (c *controller) Delete(ctx *fiber.Ctx) error {
	err := c.service.Delete(ctx.Params("id"))
	if err != nil {
		return err.ToFiber(ctx)
	}

	return ctx.SendStatus(http.StatusNoContent)
}

func (c *controller) AddComment(ctx *fiber.Ctx) error {
	var body = new(CommentRequest)
	if err := validator.BodyParser(ctx, body); err != nil {
		logger.Error(
			"error parsing body",
			"func", "AddComment - ctx.BodyParser",
			"error", err.Error(),
		)
		return errors.New(http.StatusBadRequest, err.Error(), nil).
			ToFiber(ctx)
	}

	response, err := c.service.AddComment(ctx.Locals("user").(*session.Auth), ctx.Params("id"), body)
	if err != nil {
		return err.ToFiber(ctx)
	}

	return ctx.Status(http.StatusCreated).JSON(response)
}

func (c *controller) GetByID(ctx *fiber.Ctx) error {
	response, err := c.service.GetByID(ctx.Locals("user").(*session.Auth), ctx.Params("id"))
	if err != nil {
		return err.ToFiber(ctx)
	}

	return ctx.Status(http.StatusOK).JSON(response)
}
