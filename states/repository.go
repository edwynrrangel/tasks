package states

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/edwynrrangel/tasks/config"
	"github.com/edwynrrangel/tasks/logger"
)

type repository struct {
	clientDB sqlx.DB
}

func NewRepository(config *config.Configuration) Repository {
	return &repository{
		clientDB: *config.PostgreSQLClient,
	}
}

func (r *repository) GetByID(id string) (state *StateSQL, err error) {
	state = new(StateSQL)
	query := fmt.Sprintf(sqlGetStates, " AND id = $1")

	if err := r.clientDB.Get(state, query, id); err != nil {
		logger.Error(
			"error getting state by id",
			"func", "GetByID - r.clientDB.Get",
			"error", err.Error(),
		)
		return nil, err
	}

	return
}

func (r *repository) GetByName(name string) (state *StateSQL, err error) {
	state = new(StateSQL)
	query := fmt.Sprintf(sqlGetStates, " AND name = $1")

	if err := r.clientDB.Get(state, query, name); err != nil {
		logger.Error(
			"error getting state by current status",
			"func", "GetByName - r.clientDB.Get",
			"error", err.Error(),
		)
		return nil, err
	}

	return
}

func (r *repository) GetNextStatesByID(currentStateID string) (nextStates *ListNextStatesSQL, err error) {
	statesSQL := []StateSQL{}

	query := fmt.Sprintf(sqlGetStateTransitions, " AND current_state_id = $1")
	if err := r.clientDB.Select(&statesSQL, query, currentStateID); err != nil {
		logger.Error(
			"error getting state by current status",
			"func", "GetNextStatesByID - r.clientDB.Select",
			"error", err.Error(),
		)
		return nil, err
	}

	nextStates = new(ListNextStatesSQL)
	nextStates.Status = statesSQL

	return

}
