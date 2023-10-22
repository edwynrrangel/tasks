package users

import (
	"github.com/gofiber/fiber/v2"

	"github.com/edwynrrangel/tasks/config"
	"github.com/edwynrrangel/tasks/jwt"
	"github.com/edwynrrangel/tasks/middlewares/role"
	"github.com/edwynrrangel/tasks/middlewares/session"
	"github.com/edwynrrangel/tasks/password"
	"github.com/edwynrrangel/tasks/roles"
	"github.com/edwynrrangel/tasks/token"
)

func bootstrap(config *config.Configuration) (Controller, session.Middleware, role.Middleware) {
	userRepo := NewRepository(config)
	tokenRepo := token.NewRepository(config)
	roleRepo := roles.NewRepository(config)
	jwt := jwt.NewJWT(config)
	password := password.New()
	userService := NewService(userRepo, roleRepo, password)
	tokenService := token.NewService(tokenRepo, jwt)

	return NewController(userService), session.NewSession(tokenService), role.NewRole(roleRepo)
}

func ApplyRoutes(app *fiber.App, config *config.Configuration) {
	ctrl, session, role := bootstrap(config)

	group := app.Group("/users")
	group.Post("", session.Validate, role.Authorization("Administrator"), ctrl.Create)
	group.Get("", session.Validate, role.Authorization("Administrator"), ctrl.List)
}
