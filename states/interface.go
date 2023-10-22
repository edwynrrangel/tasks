package states

type Repository interface {
	GetByID(id string) (*StateSQL, error)
	GetByName(name string) (*StateSQL, error)
	GetNextStatesByID(currentStateID string) (*ListNextStatesSQL, error)
}
