package roles

type (
	RoleSQL struct {
		ID          string `db:"id"`
		Name        string `db:"name"`
		Description string `db:"description"`
	}
)
