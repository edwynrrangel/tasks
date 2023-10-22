package auth

import (
	"github.com/edwynrrangel/tasks/errors"
	"github.com/edwynrrangel/tasks/logger"
	"github.com/edwynrrangel/tasks/middlewares/session"
	"github.com/edwynrrangel/tasks/password"
	"github.com/edwynrrangel/tasks/token"
	"github.com/edwynrrangel/tasks/users"
)

type service struct {
	userRepo     users.Repository
	tokenRepo    token.Repository
	tokenService token.Service
	password     password.Password
}

func NewService(
	userRepo users.Repository,
	tokenRepo token.Repository,
	tokenService token.Service,
	password password.Password,
) Service {
	return &service{
		userRepo:     userRepo,
		tokenRepo:    tokenRepo,
		tokenService: tokenService,
		password:     password,
	}
}

func (s *service) Login(body LoginRequest) (*LoginResponse, *errors.Error) {

	// Get user by username
	user, err := s.userRepo.GetByUsername(body.Username)
	if err != nil {
		logger.Error(
			"error getting user by username",
			"func", "Login - s.userRepo.GetByUsername",
			"error", err,
		)
		return nil, handlerError(errUnauthorized, nil)
	}

	// validate credentials
	errValidate := s.validateCredentials(body.Password, user)
	if errValidate != nil && errValidate != errFirstLogin {
		logger.Error(
			"error validating credentials",
			"func", "Login - s.validateCredentials",
			"error", errValidate,
		)
		return nil, handlerError(errValidate, nil)
	}

	// Generate token
	accessToken, err := s.generateToken(user, errValidate == errFirstLogin)
	if err != nil {
		logger.Error(
			"error generating token",
			"func", "Login - s.generateToken",
			"error", err,
		)

		return nil, handlerError(err, nil)
	}

	if errValidate == errFirstLogin {
		details := errors.Details{
			"access_token": accessToken,
		}
		return nil, handlerError(errValidate, details)
	}

	return &LoginResponse{accessToken}, nil
}

func (s *service) ChangePassword(authUser *session.Auth, body ChangePasswordRequest) (*LoginResponse, *errors.Error) {
	user, err := s.userRepo.GetByID(authUser.ID)
	if err != nil {
		logger.Error(
			"error getting user by username",
			"func", "Login - s.userRepo.GetByUsername",
			"error", err,
		)
		return nil, handlerError(errUnauthorized, nil)
	}

	if err := s.validateCredentials(body.CurrentPassword, user); err != nil &&
		err != errFirstLogin {
		logger.Error(
			"error validating credentials",
			"func", "Login - s.validateCredentials",
			"error", err,
		)
		return nil, handlerError(err, nil)
	}

	password, err := s.password.Hash(body.NewPassword)
	if err != nil {
		logger.Error(
			"error hashing password",
			"func", "Login - hashPassword",
			"error", err,
		)
		return nil, handlerError(err, nil)
	}

	if err := s.userRepo.UpdatePassword(user.ID, password); err != nil {
		logger.Error(
			"error updating user password",
			"func", "Login - s.userRepo.UpdatePassword",
			"error", err,
		)
		return nil, handlerError(err, nil)
	}
	// generate token
	accessToken, err := s.generateToken(user, false)
	if err != nil {
		logger.Error(
			"error generating token",
			"func", "Login - s.generateToken",
			"error", err,
		)
		return nil, handlerError(err, nil)
	}
	// return token
	return &LoginResponse{accessToken}, nil
}

func (s *service) Logout(authUser *session.Auth) *errors.Error {
	if err := s.tokenRepo.Delete(authUser.ID); err != nil {
		logger.Error(
			"error deleting token",
			"func", "Logout - s.tokenRepo.Delete",
			"error", err,
		)
		return handlerError(err, nil)
	}

	return nil
}

func (s *service) validateCredentials(password string, user *users.UserSQL) error {

	// validate if is first login, password is empty and role is Administrator return errFirstLogin
	if user.FirstLogin && !user.Password.Valid && user.Role.Name == "Administrator" {
		return errFirstLogin
	}
	// validate password if role is not Administrator but is first login
	if user.FirstLogin && user.Role.Name != "Administrator" {
		if !s.password.IsValid(user.Password.String, password) {
			return errUnauthorized
		}
		return errFirstLogin
	}

	if !s.password.IsValid(user.Password.String, password) {
		return errUnauthorized
	}

	if !user.Active {
		return errInactiveUser
	}

	return nil
}

func (s *service) generateToken(user *users.UserSQL, firstLogin bool) (string, error) {
	claim := struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Role     string `json:"role"`
	}{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role.Name,
	}
	token, err := s.tokenService.Generate(claim, firstLogin)
	if err != nil {
		logger.Error(
			"error generating token",
			"func", "Login - s.tokenUtil.GenerateToken",
			"error", err,
		)
		return "", err
	}

	return token, nil
}
