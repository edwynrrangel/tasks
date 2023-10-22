package tasks

import (
	"github.com/gofiber/fiber/v2"

	"github.com/edwynrrangel/tasks/comments"
	"github.com/edwynrrangel/tasks/config"
	"github.com/edwynrrangel/tasks/jwt"
	"github.com/edwynrrangel/tasks/middlewares/session"
	"github.com/edwynrrangel/tasks/states"
	"github.com/edwynrrangel/tasks/token"
	"github.com/edwynrrangel/tasks/users"
)

func bootstrap(config *config.Configuration) (Controller, session.Middleware) {
	taskRepo := NewRepository(config)
	userRepo := users.NewRepository(config)
	stateRepo := states.NewRepository(config)
	commentRepo := comments.NewRepository(config)
	tokenRepo := token.NewRepository(config)
	jwt := jwt.NewJWT(config)
	taskService := NewService(taskRepo, userRepo, stateRepo, commentRepo)
	tokenService := token.NewService(tokenRepo, jwt)

	return NewController(taskService), session.NewSession(tokenService)
}

func ApplyRoutes(app *fiber.App, config *config.Configuration) {
	ctrl, middleware := bootstrap(config)

	group := app.Group("/tasks")
	group.Post("", middleware.ValidateSession, ctrl.Create)
	group.Get("", middleware.ValidateSession, ctrl.List)
	group.Patch("/:id", middleware.ValidateSession, ctrl.Update)
	group.Delete("/:id", middleware.ValidateSession, ctrl.Delete)
	group.Post("/:id/comments", middleware.ValidateSession, ctrl.AddComment)
	group.Get("/:id", middleware.ValidateSession, ctrl.GetByID)
}
