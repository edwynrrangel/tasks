package roles

type Repository interface {
	GetByID(id string) (*RoleSQL, error)
	GetByName(name string) (*RoleSQL, error)
}
