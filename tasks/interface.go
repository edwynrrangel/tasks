package tasks

import (
	"github.com/gofiber/fiber/v2"

	"github.com/edwynrrangel/tasks/comments"
	"github.com/edwynrrangel/tasks/errors"
	"github.com/edwynrrangel/tasks/middlewares/session"
)

type Repository interface {
	Create(task *TaskSQL) error
	List(queryParems ListRequest) ([]TaskSQL, uint, error)
	GetByID(id string) (*TaskSQL, error)
	Update(id string, task *TaskSQL) error
	Delete(id string) error
}

type Service interface {
	Create(task *CreateRequest) (*TaskResponse, *errors.Error)
	List(queryParems ListRequest) (*ListResponse, *errors.Error)
	Update(id string, task *UpdateRequest) (*TaskResponse, *errors.Error)
	UpdateStatus(id string, body *UpdateRequest) (*TaskResponse, *errors.Error)
	Delete(id string) *errors.Error
	AddComment(user *session.Auth, id string, body *CommentRequest) (*comments.CommentResponse, *errors.Error)
	GetByID(id string) (*TaskResponse, *errors.Error)
}

type Controller interface {
	Create(ctx *fiber.Ctx) error
	List(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	AddComment(ctx *fiber.Ctx) error
	GetByID(ctx *fiber.Ctx) error
}
