package tasks

import (
	"errors"
)

const (
	errMsgRoleNotExecutor   string = "el usuario tiene que tener el rol de Executor"
	errMsgTaskNotFound      string = "no se encontró la tarea"
	errMsgStatusNotStarted  string = "el estado de la tarea no es Asignado"
	errMsgStatusNotAllow    string = "el nuevo estado no es permitido"
	errMsgDueDateExpired    string = "la tarea ya expiró"
	errMsgDueDateNotExpired string = "la tarea no ha expirado"
	errMsgNotAuthorized     string = "no tiene permiso para realizar esta acción"
)
const (
	sqlInsert = `
		INSERT INTO task.tasks 
			(title, description, due_date, status, assigned_user) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	sqlGetTask = `
		SELECT
			t.id
			,t.title
			,t.description
			,t.due_date
			,u.id AS "assigned_user.id" 
			,u.username AS "assigned_user.username" 
			,t.created_at
			,t.updated_at
			,s.id AS "status.id"
			,s.name AS "status.name"
		FROM task.tasks t
		INNER JOIN usr.users u
			ON t.assigned_user = u.id
		INNER JOIN task.states s
			ON t.status = s.id
		WHERE 1 = 1 %s
	`

	sqlCountTasks = `
	SELECT COUNT(*) FROM task.tasks t
	INNER JOIN usr.users u
		ON t.assigned_user = u.id
	INNER JOIN task.states s
		ON t.status = s.id
	WHERE 1 = 1 %s
	`
)

var (
	errNotFound error = errors.New(errMsgTaskNotFound)
)
