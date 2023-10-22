package users

import "errors"

const (
	errMsgRoleInvalid       string = "rol no permitido"
	errMgsUserAlreadyExists string = "el usuario ya existe"
)

const (
	sqlGetUser = `
	SELECT 
		u.id
		,u.username
		,u.password
		,u.first_login
		,u.created_at
		,u.updated_at
		,u.last_login_at
		,r.id AS "role.id"
		,r.name AS "role.name"
		,r.description AS "role.description"
		,u.active
	FROM usr.users u
	INNER JOIN usr.roles r ON r.id = u.role_id
	WHERE 1 = 1 %s
	`

	sqlUpdatePassword = `
	UPDATE usr.users
	SET 
		password = $1
		,first_login = false
		,updated_at = NOW()
	WHERE id = $2
	`

	sqlSave = `
	INSERT INTO usr.users (
		username
		,password
		,role_id
	)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at
	`
	sqlCountUsers = `
	SELECT COUNT(*) FROM usr.users u
	INNER JOIN usr.roles r ON r.id = u.role_id
	WHERE 1 = 1 %s
	`
)

var (
	errUserAlreadyExists error = errors.New(errMgsUserAlreadyExists)
)
