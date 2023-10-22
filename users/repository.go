package users

import (
	"fmt"
	"strings"

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

func (r *repository) GetByUsername(username string) (user *UserSQL, err error) {
	user = new(UserSQL)
	query := fmt.Sprintf(sqlGetUser, "AND username = $1")

	if err = r.clientDB.Get(user, query, username); err != nil {
		logger.Error(
			"error getting user by username",
			"func", "GetByUsername",
			"error", err.Error(),
		)
		return nil, err
	}
	return
}

func (r *repository) GetByID(id string) (user *UserSQL, err error) {
	user = new(UserSQL)
	query := fmt.Sprintf(sqlGetUser, "AND u.id = $1")

	if err = r.clientDB.Get(user, query, id); err != nil {
		logger.Error(
			"error getting user by id",
			"func", "GetByID",
			"error", err.Error(),
		)
		return nil, err
	}
	return
}

func (r *repository) UpdatePassword(id, password string) (err error) {
	if _, err = r.clientDB.Exec(sqlUpdatePassword, password, id); err != nil {
		logger.Error(
			"error updating user password",
			"func", "UpdatePassword",
			"error", err.Error(),
		)
		return err
	}
	return
}

func (r *repository) Create(user *UserSQL) error {
	err := r.clientDB.QueryRow(
		sqlSave,
		user.Username,
		user.Password,
		user.Role.ID,
	).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		logger.Error(
			"error saving user",
			"func", "Create - r.clientDB.Exec",
			"error", err.Error(),
		)
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return errUserAlreadyExists
		}
		return err
	}

	return nil
}

func (r *repository) List(queryParams ListRequest) (users []UserSQL, total uint, err error) {
	fl := new(filterList)
	fl.getUsername(queryParams).
		getRole(queryParams).
		getCreatedAt(queryParams)

	query := fmt.Sprintf(sqlCountUsers, fl.filters)
	if err = r.clientDB.Get(&total, query, fl.args...); err != nil {
		logger.Error(
			"error getting total users",
			"func", "List - r.clientDB.Get",
			"error", err.Error(),
		)
		return nil, total, err
	}

	if total == 0 {
		return
	}

	fl.getOrderBy(queryParams).
		getLimit(queryParams).
		getOffset(queryParams)

	query = fmt.Sprintf(sqlGetUser, fl.filters)
	if err = r.clientDB.Select(&users, query, fl.args...); err != nil {
		logger.Error(
			"error getting users",
			"func", "List - r.clientDB.Select",
			"error", err.Error(),
		)
		return nil, total, err
	}

	return
}
