package token

import (
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

func (r *repository) GetByUserID(userID string) (*TokenSQL, error) {
	tokenSQL := new(TokenSQL)
	query := `SELECT * FROM usr.tokens WHERE user_id = $1`
	if err := r.clientDB.Get(tokenSQL, query, userID); err != nil {
		logger.Error(
			"error getting token by user id",
			"func", "GetByUserID - r.clientDB.Get",
			"error", err.Error(),
		)
		return nil, err
	}

	return tokenSQL, nil
}

func (r *repository) GetByToken(token string) (*TokenSQL, error) {
	tokenSQL := new(TokenSQL)
	query := `SELECT * FROM usr.tokens WHERE jwt_token = $1`
	if err := r.clientDB.Get(tokenSQL, query, token); err != nil {
		logger.Error(
			"error getting token by token",
			"func", "GetByToken - r.clientDB.Get",
			"error", err.Error(),
		)
		return nil, err
	}

	return tokenSQL, nil
}

func (r *repository) Save(token *TokenSQL) error {
	query := `
		INSERT INTO usr.tokens (user_id, jwt_token, expiration, first_login) 
		VALUES ($1, $2, $3, $4)
	`
	if _, err := r.clientDB.Exec(query, token.UserID, token.JWTToken, token.Expiration, token.FirstLogin); err != nil {
		logger.Error(
			"error saving token",
			"func", "Save - r.clientDB.Exec",
			"error", err.Error(),
		)
		return err
	}

	return nil
}

func (r *repository) Update(userID string, token *TokenSQL) error {
	query := `
		UPDATE usr.tokens 
		SET jwt_token = $1, expiration = $2, first_login = $3
		WHERE user_id = $4
	`
	result, err := r.clientDB.Exec(query, token.JWTToken, token.Expiration, token.FirstLogin, token.UserID)
	if err != nil {
		logger.Error(
			"error updating token",
			"func", "Update - r.clientDB.Exec",
			"error", err.Error(),
		)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error(
			"error getting rows affected",
			"func", "Update - result.RowsAffected",
			"error", err.Error(),
		)
		return err
	}

	if rowsAffected == 0 {
		return errNotFound
	}

	return nil
}

func (r *repository) UserIsFirstLogin(userID string) (bool, error) {
	userSQL := new(UserSQL)
	query := `SELECT * FROM usr.users WHERE id = $1`
	if err := r.clientDB.Get(userSQL, query, userID); err != nil {
		logger.Error(
			"error getting user by id",
			"func", "UserIsFirstLogin - r.clientDB.Get",
			"error", err.Error(),
		)
		return false, err
	}

	return userSQL.FirstLogin, nil
}

func (r *repository) Delete(userID string) error {
	query := `
		DELETE FROM usr.tokens
		WHERE user_id = $1
	`

	result, err := r.clientDB.Exec(query, userID)
	if err != nil {
		logger.Error(
			"error deleting token",
			"func", "Delete - r.clientDB.Exec",
			"error", err.Error(),
		)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error(
			"error getting rows affected",
			"func", "Delete - result.RowsAffected",
			"error", err.Error(),
		)
		return err
	}

	if rowsAffected == 0 {
		return errNotFound
	}

	return nil
}
