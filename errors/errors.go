package errors

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type Details map[string]interface{}

type Error struct {
	Code    int                    `json:"-"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func New(code int, message string, datails Details) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Details: datails,
	}
}

func (e *Error) ToFiber(ctx *fiber.Ctx) error {
	return ctx.Status(e.Code).JSON(e)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}
