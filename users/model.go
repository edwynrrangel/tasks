package users

import (
	"database/sql"
)

type (
	RoleSQL struct {
		ID          string `db:"id"`
		Name        string `db:"name"`
		Description string `db:"description"`
	}

	UserSQL struct {
		ID          string         `db:"id"`
		Username    string         `db:"username"`
		Password    sql.NullString `db:"password"`
		Role        RoleSQL        `db:"role"`
		FirstLogin  bool           `db:"first_login"`
		CreatedAt   string         `db:"created_at"`
		UpdatedAt   string         `db:"updated_at"`
		LastLoginAt sql.NullTime   `db:"last_login_at"`
		Active      bool           `db:"active"`
	}
)

func (u *UserSQL) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Role:      u.Role.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
