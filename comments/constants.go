package comments

const (
	sqlInsert = `
	INSERT INTO task.comments
		(task_id, comment, created_by)
	VALUES ($1, $2, $3)
	RETURNING id, created_at
	`

	sqlGet = `
	SELECT
		c.id
		,c.comment
		,c.created_at
		,u.id AS "created_by.id"
		,u.username AS "created_by.username"
	FROM task.comments c
	INNER JOIN usr.users u 
		ON u.id = c.created_by
	WHERE 1 = 1 %s
	`
)
