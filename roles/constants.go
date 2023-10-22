package roles

const (
	sqlGetRole string = `
	SELECT 
		* 
	FROM usr.roles 
	WHERE 1 = 1 %s
	`
)
