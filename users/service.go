package users

import (
	"database/sql"
	"net/http"

	"github.com/edwynrrangel/tasks/errors"
	"github.com/edwynrrangel/tasks/logger"
	"github.com/edwynrrangel/tasks/password"
	"github.com/edwynrrangel/tasks/roles"
)

type service struct {
	userRepo Repository
	roleRepo roles.Repository
	password password.Password
}

func NewService(
	userRepo Repository,
	roleRepo roles.Repository,
	password password.Password,
) Service {
	return &service{
		userRepo: userRepo,
		roleRepo: roleRepo,
		password: password,
	}
}

func (s *service) Create(user *CreateRequest) (*UserResponse, *errors.Error) {

	if user.Role == "Administrator" {
		return nil, errors.New(http.StatusConflict, errMsgRoleInvalid, nil)
	}

	// Generate random password
	password, err := s.password.GenerateRamdom(user.PasswordLength)
	if err != nil {
		logger.Error(
			"error generating random password",
			"func", "generateRandomPassword",
			"error", err.Error(),
		)
		return nil, handlerError(err, nil)
	}

	// Get role
	role, err := s.roleRepo.GetByName(user.Role)
	if err != nil {
		logger.Error(
			"error getting role by name",
			"func", "GetByName",
			"error", err.Error(),
		)
		return nil, handlerError(err, nil)
	}

	// hash password
	passwordHash, err := s.password.Hash(password)
	if err != nil {
		logger.Error(
			"error hashing password",
			"func", "Hash",
			"error", err.Error(),
		)
		return nil, handlerError(err, nil)
	}

	userSQL := &UserSQL{
		Username: user.Username,
		Password: sql.NullString{
			String: passwordHash,
			Valid:  true,
		},
		Role: (RoleSQL)(*role),
	}

	if err := s.userRepo.Create(userSQL); err != nil {
		logger.Error(
			"error saving user",
			"func", "Save",
			"error", err.Error(),
		)
		return nil, handlerError(err, nil)
	}

	userResponse := userSQL.ToResponse()
	userResponse.Password = password
	return userResponse, nil
}

func (s *service) List(queryParams ListRequest) (*ListResponse, *errors.Error) {
	list, total, err := s.userRepo.List(queryParams)
	if err != nil {
		logger.Error(
			"error getting users",
			"func", "List - s.userRepo.List",
			"error", err.Error(),
		)
		return nil, handlerError(err, nil)
	}

	users := make([]UserResponse, 0, len(list))
	for _, user := range list {
		users = append(users, *user.ToResponse())
	}

	return &ListResponse{
		Total: total,
		Users: users,
	}, nil
}
