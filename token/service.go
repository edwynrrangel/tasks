package token

import (
	"database/sql"
	"net/http"

	"github.com/edwynrrangel/tasks/errors"
	"github.com/edwynrrangel/tasks/jwt"
	"github.com/edwynrrangel/tasks/logger"
)

type service struct {
	tokenRepo Repository
	jwt       jwt.JWT
}

func NewService(tokenRepo Repository, jwt jwt.JWT) Service {
	return &service{
		tokenRepo: tokenRepo,
		jwt:       jwt,
	}
}

func (s *service) Generate(payload interface{}, firstLogin bool) (string, error) {

	// Sign and return token
	tokenString, expiresAt, err := s.jwt.Generate(payload)
	if err != nil {
		logger.Error(
			"error signing token",
			"func", "GenerateToken - s.jwt.Generate",
			"error", err,
		)
		return "", err
	}

	// Get token by user id
	currentClaim := payload.(struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Role     string `json:"role"`
	})

	tokenSQL, err := s.tokenRepo.GetByUserID(currentClaim.ID)
	if err != nil && err != sql.ErrNoRows {
		logger.Error(
			"error getting token by user id",
			"func", "GenerateToken - s.tokenRepo.GetByUserID",
			"error", err,
		)
		return "", err
	}

	// Update token
	if tokenSQL != nil {
		tokenSQL.JWTToken = tokenString
		tokenSQL.Expiration = expiresAt
		tokenSQL.FirstLogin = firstLogin

		if err := s.tokenRepo.Update(currentClaim.ID, tokenSQL); err != nil {
			logger.Error(
				"error updating token",
				"func", "GenerateToken - s.tokenRepo.Update",
				"error", err,
			)
			return "", err
		}
	}

	// Save token
	if tokenSQL == nil {
		tokenSQL = &TokenSQL{
			UserID:     currentClaim.ID,
			JWTToken:   tokenString,
			Expiration: expiresAt,
			FirstLogin: firstLogin,
		}

		if err := s.tokenRepo.Save(tokenSQL); err != nil {
			logger.Error(
				"error saving token",
				"func", "GenerateToken - s.tokenRepo.Save",
				"error", err,
			)
			return "", err
		}
	}

	return tokenString, nil
}

func (s *service) Validate(token string) *errors.Error {
	// Check if token belongs to user
	tokenSQL, err := s.tokenRepo.GetByToken(token)
	if err != nil {
		logger.Error(
			"error getting token by token",
			"func", "ValidateToken - s.tokenRepo.GetByToken",
			"error", err,
		)
		return errors.New(http.StatusUnauthorized, errMsgExpiredInvalidMissing, nil)
	}

	// Check if token is expired, invalid or missing
	if !s.jwt.IsValid(tokenSQL.JWTToken) {
		return errors.New(http.StatusUnauthorized, errMsgExpiredInvalidMissing, nil)
	}

	// Check if token was generated for first login
	if tokenSQL.FirstLogin {
		return errors.New(http.StatusUnauthorized, errMsgFirstLogin, nil)
	}

	return nil
}

func (s *service) ValidateFirstLogin(token, userID string) *errors.Error {
	// Check if token belongs to user
	tokenSQL, err := s.tokenRepo.GetByToken(token)
	if err != nil {
		logger.Error(
			"error getting token by token",
			"func", "ValidateFirstLogin - s.tokenRepo.GetByToken",
			"error", err,
		)
		return errors.New(http.StatusUnauthorized, errMsgExpiredInvalidMissing, nil)
	}

	// Check if token is expired, invalid or missing
	if !s.jwt.IsValid(tokenSQL.JWTToken) {
		return errors.New(http.StatusUnauthorized, errMsgExpiredInvalidMissing, nil)
	}

	return nil
}
