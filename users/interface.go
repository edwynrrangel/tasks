package users

import (
	"github.com/gofiber/fiber/v2"

	"github.com/edwynrrangel/tasks/errors"
)

type Repository interface {
	GetByUsername(username string) (*UserSQL, error)
	GetByID(id string) (*UserSQL, error)
	UpdatePassword(id, password string) error
	Create(user *UserSQL) error
	List(queryParams ListRequest) ([]UserSQL, uint, error)
}

type Service interface {
	Create(user *CreateRequest) (*UserResponse, *errors.Error)
	List(queryParams ListRequest) (*ListResponse, *errors.Error)
}

type Controller interface {
	Create(ctx *fiber.Ctx) error
	List(ctx *fiber.Ctx) error
}
