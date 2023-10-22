package token

import (
	"time"
)

type (
	TokenSQL struct {
		UserID     string    `db:"user_id"`
		JWTToken   string    `db:"jwt_token"`
		Expiration time.Time `db:"expiration"`
		CreatedAt  time.Time `db:"created_at"`
		FirstLogin bool      `db:"first_login"`
	}

	UserSQL struct {
		ID         string `db:"id"`
		FirstLogin bool   `db:"first_login"`
	}
)
