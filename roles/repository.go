package roles

import (
	"fmt"

	"github.com/edwynrrangel/tasks/config"
	"github.com/edwynrrangel/tasks/logger"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	clientDB *sqlx.DB
}

func NewRepository(config *config.Configuration) Repository {
	return &repository{
		clientDB: config.PostgreSQLClient,
	}
}

func (r *repository) GetByID(id string) (*RoleSQL, error) {
	var roleSQL = new(RoleSQL)
	query := fmt.Sprintf(sqlGetRole, " AND id = $1")
	err := r.clientDB.Get(roleSQL, query, id)
	if err != nil {
		logger.Error(
			"error getting role by id",
			"func", "GetByID - r.clientDB.Get",
			"error", err.Error(),
		)
		return nil, err
	}
	return roleSQL, nil
}

func (r *repository) GetByName(name string) (*RoleSQL, error) {
	var roleSQL = new(RoleSQL)
	query := fmt.Sprintf(sqlGetRole, " AND name = $1")
	err := r.clientDB.Get(roleSQL, query, name)
	if err != nil {
		logger.Error(
			"error getting role by name",
			"func", "GetByName - r.clientDB.Get",
			"error", err.Error(),
		)
		return nil, err
	}
	return roleSQL, nil
}
