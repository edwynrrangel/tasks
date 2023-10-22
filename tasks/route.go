package tasks

import (
	"github.com/gofiber/fiber/v2"

	"github.com/edwynrrangel/tasks/comments"
	"github.com/edwynrrangel/tasks/config"
	"github.com/edwynrrangel/tasks/jwt"
	"github.com/edwynrrangel/tasks/middlewares/role"
	"github.com/edwynrrangel/tasks/middlewares/session"
	"github.com/edwynrrangel/tasks/roles"
	"github.com/edwynrrangel/tasks/states"
	"github.com/edwynrrangel/tasks/token"
	"github.com/edwynrrangel/tasks/users"
)

func bootstrap(config *config.Configuration) (Controller, session.Middleware, role.Middleware) {
	taskRepo := NewRepository(config)
	userRepo := users.NewRepository(config)
	stateRepo := states.NewRepository(config)
	commentRepo := comments.NewRepository(config)
	tokenRepo := token.NewRepository(config)
	roleRepo := roles.NewRepository(config)
	jwt := jwt.NewJWT(config)
	taskService := NewService(taskRepo, userRepo, stateRepo, commentRepo)
	tokenService := token.NewService(tokenRepo, jwt)

	return NewController(taskService), session.NewSession(tokenService), role.NewRole(roleRepo)
}

func ApplyRoutes(app *fiber.App, config *config.Configuration) {
	ctrl, session, role := bootstrap(config)

	group := app.Group("/tasks")
	group.Post("", session.Validate, role.Authorization("Administrator"), ctrl.Create)
	group.Get("", session.Validate, ctrl.List)
	group.Patch("/:id", session.Validate, role.Authorization("Administrator", "Executor"), ctrl.Update)
	group.Delete("/:id", session.Validate, role.Authorization("Administrator"), ctrl.Delete)
	group.Post("/:id/comments", session.Validate, role.Authorization("Executor"), ctrl.AddComment)
	group.Get("/:id", session.Validate, ctrl.GetByID)
}
