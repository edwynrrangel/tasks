package token

import "github.com/edwynrrangel/tasks/errors"

type Repository interface {
	GetByUserID(userID string) (*TokenSQL, error)
	GetByToken(token string) (*TokenSQL, error)
	Save(token *TokenSQL) error
	Update(userID string, token *TokenSQL) error
	UserIsFirstLogin(userID string) (bool, error)
	Delete(token string) error
}

type Service interface {
	Generate(payload interface{}, firstLogin bool) (string, error)
	Validate(token string) *errors.Error
	ValidateFirstLogin(token, userID string) *errors.Error
}
