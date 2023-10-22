package tasks

import (
	"net/http"

	"github.com/edwynrrangel/tasks/comments"
	"github.com/edwynrrangel/tasks/errors"
	"github.com/edwynrrangel/tasks/logger"
	"github.com/edwynrrangel/tasks/middlewares/session"
	"github.com/edwynrrangel/tasks/states"
	"github.com/edwynrrangel/tasks/users"
)

type service struct {
	taskRepo    Repository
	userRepo    users.Repository
	stateRepo   states.Repository
	commentRepo comments.Repository
}

func NewService(
	taskRepo Repository,
	userRepo users.Repository,
	stateRepo states.Repository,
	commentRepo comments.Repository,
) Service {
	return &service{
		taskRepo:    taskRepo,
		userRepo:    userRepo,
		stateRepo:   stateRepo,
		commentRepo: commentRepo,
	}
}

func (s *service) Create(task *CreateRequest) (*TaskResponse, *errors.Error) {
	user, err := s.userRepo.GetByID(task.AssignedUser)
	if err != nil {
		logger.Error(
			"error getting user by id",
			"func", "Create - s.userRepo.GetByID",
			"error", err.Error(),
		)
		return nil, handlerError(err, nil)
	}

	if user.Role.Name != "Executor" {
		return nil, errors.New(http.StatusConflict, errMsgRoleNotExecutor, nil)
	}

	state, err := s.stateRepo.GetByName(task.Status)
	if err != nil {
		logger.Error(
			"error getting state by id",
			"func", "Create - s.stateRepo.GetByID",
			"error", err.Error(),
		)
		return nil, handlerError(err, nil)
	}

	taskSQL := task.ToTaskSQL(state.ID)
	err = s.taskRepo.Create(taskSQL)
	if err != nil {
		logger.Error(
			"error creating task",
			"func", "Create - s.taskRepo.Create",
			"error", err.Error(),
		)
		return nil, handlerError(err, nil)
	}

	taskResponse := taskSQL.ToResponse()
	taskResponse.Status.Name = state.Name
	taskResponse.AssignedUser.Username = user.Username
	return taskResponse, nil
}

func (s *service) List(user *session.Auth, queryParems ListRequest) (*ListResponse, *errors.Error) {
	if user.Role == "Executor" {
		queryParems.Username = user.Username
	}

	list, total, err := s.taskRepo.List(queryParems)
	if err != nil {
		logger.Error(
			"error getting tasks",
			"func", "List - s.taskRepo.List",
			"error", err.Error(),
		)
		return nil, handlerError(err, nil)
	}

	tasks := make([]TaskResponse, 0, len(list))
	for _, task := range list {
		tasks = append(tasks, *task.ToResponse())
	}

	return &ListResponse{
		Total: total,
		Tasks: tasks,
	}, nil
}

func (s *service) Update(id string, body *UpdateRequest) (*TaskResponse, *errors.Error) {
	taskSQL, err := s.taskRepo.GetByID(id)
	if err != nil {
		logger.Error(
			"error getting task by id",
			"func", "Update - s.taskRepo.GetByID",
			"error", err.Error(),
		)
		return nil, handlerError(err, nil)
	}

	if taskSQL.Status.Name != "Asignado" {
		return nil, errors.New(http.StatusConflict, errMsgStatusNotStarted, nil)
	}

	user, err := s.userRepo.GetByID(body.AssignedUser)
	if err != nil {
		logger.Error(
			"error getting user by id",
			"func", "Update - s.userRepo.GetByID",
			"error", err.Error(),
		)
		return nil, handlerError(err, nil)
	}

	if user.Role.Name != "Executor" {
		return nil, errors.New(http.StatusConflict, errMsgRoleNotExecutor, nil)
	}

	taskSQL.
		updateTitle(body).
		updateDescription(body).
		updateDueDate(body).
		updateAssignedUser(body)

	err = s.taskRepo.Update(id, taskSQL)
	if err != nil {
		logger.Error(
			"error updating task",
			"func", "Update - s.taskRepo.Update",
			"error", err.Error(),
		)
		return nil, handlerError(err, nil)
	}

	taskResponse := taskSQL.ToResponse()
	taskResponse.AssignedUser.Username = user.Username

	return taskResponse, nil
}

func (s *service) UpdateStatus(id string, body *UpdateRequest) (*TaskResponse, *errors.Error) {
	taskSQL, err := s.taskRepo.GetByID(id)
	if err != nil {
		logger.Error(
			"error getting task by id",
			"func", "Update - s.taskRepo.GetByID",
			"error", err.Error(),
		)
		return nil, handlerError(err, nil)
	}

	// validate is DueDate expired
	if taskSQL.IsExpired() {
		return nil, errors.New(http.StatusConflict, errMsgDueDateExpired, nil)
	}

	listNextStatusSQL, err := s.stateRepo.GetNextStatesByID(taskSQL.Status.ID)
	if err != nil {
		logger.Error(
			"error getting next status",
			"func", "UpdateStatus - s.stateRepo.GetNextStatus",
			"error", err.Error(),
		)
		return nil, handlerError(err, nil)
	}

	// validate if the status is valid
	nextStatus := listNextStatusSQL.IsValidNextStateByName(body.Status)
	if nextStatus == nil {
		return nil, errors.New(http.StatusConflict, errMsgStatusNotAllow, nil)
	}

	taskSQL.updateStatus((StatusSQL)(*nextStatus))

	err = s.taskRepo.Update(id, taskSQL)
	if err != nil {
		logger.Error(
			"error updating task",
			"func", "Update - s.taskRepo.Update",
			"error", err.Error(),
		)
		return nil, handlerError(err, nil)
	}

	taskResponse := taskSQL.ToResponse()
	taskResponse.Status.Name = nextStatus.Name

	return taskResponse, nil
}

func (s *service) Delete(id string) *errors.Error {
	taskSQL, err := s.taskRepo.GetByID(id)
	if err != nil {
		logger.Error(
			"error getting task by id",
			"func", "Delete - s.taskRepo.GetByID",
			"error", err.Error(),
		)
		return handlerError(err, nil)
	}

	if taskSQL.Status.Name != "Asignado" {
		return errors.New(http.StatusConflict, errMsgStatusNotStarted, nil)
	}

	err = s.taskRepo.Delete(id)
	if err != nil {
		logger.Error(
			"error deleting task",
			"func", "Delete - s.taskRepo.Delete",
			"error", err.Error(),
		)
		return handlerError(err, nil)
	}

	return nil
}

func (s *service) AddComment(user *session.Auth, id string, body *CommentRequest) (*comments.CommentResponse, *errors.Error) {
	taskSQL, err := s.taskRepo.GetByID(id)
	if err != nil {
		logger.Error(
			"error getting task by id",
			"func", "AddComment - s.taskRepo.GetByID",
			"error", err.Error(),
		)
		return nil, handlerError(err, nil)
	}

	commentSQL := comments.CommentSQL{
		Task: &comments.TaskSQL{
			ID:    taskSQL.ID,
			Title: taskSQL.Title,
		},
		Comment: body.Comment,
		CreatedBy: comments.UserSQL{
			ID:       user.ID,
			Username: user.Username,
		},
	}

	err = s.commentRepo.Add(&commentSQL)
	if err != nil {
		logger.Error(
			"error adding comment",
			"func", "AddComment - s.commentRepo.Add",
			"error", err.Error(),
		)
		return nil, handlerError(err, nil)
	}

	return commentSQL.ToResponse(), nil
}

func (s *service) GetByID(user *session.Auth, id string) (*TaskResponse, *errors.Error) {
	taskSQL, err := s.taskRepo.GetByID(id)
	if err != nil {
		logger.Error(
			"error getting task by id",
			"func", "GetByID - s.taskRepo.GetByID",
			"error", err.Error(),
		)
		return nil, handlerError(err, nil)
	}

	if user.Role == "Executor" && taskSQL.AssignedUser.ID != user.ID {
		return nil, errors.New(http.StatusForbidden, errMsgNotAuthorized, nil)
	}

	taskResponse := taskSQL.ToResponse()

	commentsSQL, err := s.commentRepo.GetByTaskID(id)
	if err != nil {
		logger.Error(
			"error getting comments by task id",
			"func", "GetByID - s.commentRepo.GetByTaskID",
			"error", err.Error(),
		)
		return nil, handlerError(err, nil)
	}

	comments := make([]comments.CommentResponse, 0, len(commentsSQL))
	for _, comment := range commentsSQL {
		commentResponse := comment.ToResponse()
		comments = append(comments, *commentResponse)
	}

	taskResponse.Comments = comments
	return taskResponse, nil
}
