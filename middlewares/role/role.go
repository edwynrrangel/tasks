package role

import (
	"net/http"

	"github.com/edwynrrangel/tasks/middlewares/session"
	"github.com/edwynrrangel/tasks/roles"
	"github.com/gofiber/fiber/v2"
)

type role struct {
	repoRole roles.Repository
}

func NewRole(repoRole roles.Repository) Middleware {
	return &role{
		repoRole: repoRole,
	}
}

func (r *role) Authorization(allowedRoles ...string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Get user  from JWT
		authUser := ctx.Locals("user").(*session.Auth)
		role, err := r.repoRole.GetByName(authUser.Role)
		if err != nil {
			return ctx.Status(http.StatusForbidden).JSON(fiber.Map{
				"message": errMsgForbidden,
			})
		}

		// Check if user role is allowed
		for _, allowedRole := range allowedRoles {
			if role.Name == allowedRole {
				return ctx.Next()
			}
		}
		return ctx.Status(http.StatusForbidden).JSON(fiber.Map{
			"message": errMsgForbidden,
		})
	}
}
