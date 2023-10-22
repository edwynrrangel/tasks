package states

type (
	StateSQL struct {
		ID   string `db:"id"`
		Name string `db:"name"`
	}

	ListNextStatesSQL struct {
		Status []StateSQL
	}
)

func (l *ListNextStatesSQL) IsValidNextStateByName(name string) *StateSQL {
	for _, state := range l.Status {
		if state.Name == name {
			return &state
		}
	}

	return nil
}
