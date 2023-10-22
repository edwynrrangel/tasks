package auth

import (
	"github.com/gofiber/fiber/v2"

	"github.com/edwynrrangel/tasks/config"
	"github.com/edwynrrangel/tasks/jwt"
	"github.com/edwynrrangel/tasks/middlewares/session"
	"github.com/edwynrrangel/tasks/password"
	"github.com/edwynrrangel/tasks/token"
	"github.com/edwynrrangel/tasks/users"
)

func bootstrap(config *config.Configuration) (Controller, session.Middleware) {
	userRepo := users.NewRepository(config)
	tokenRepo := token.NewRepository(config)
	jwt := jwt.NewJWT(config)
	password := password.New()
	tokenService := token.NewService(tokenRepo, jwt)
	cs := NewService(userRepo, tokenRepo, tokenService, password)

	return NewController(cs), session.NewSession(tokenService)
}

func ApplyRoutes(app *fiber.App, config *config.Configuration) {
	ctrl, session := bootstrap(config)

	group := app.Group("/auth")
	group.Post("/login", ctrl.Login)
	group.Post("/change-password", session.ValidateFirstLogin, ctrl.ChangePassword)
	group.Post("/logout", session.ValidateFirstLogin, ctrl.Logout)
}
