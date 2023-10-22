package tasks

import (
	"fmt"

	"github.com/edwynrrangel/tasks/comments"
)

type (
	CreateRequest struct {
		Title        string `json:"title" validate:"required,min=3,max=100"`
		Description  string `json:"description" validate:"required,min=3,max=1000"`
		DueDate      string `json:"due_date" validate:"required,datetime=2006-01-02"`
		Status       string `json:"status" validate:"required,oneof=Asignado"`
		AssignedUser string `json:"assigned_user" validate:"required,uuid4"`
	}

	ListRequest struct {
		Page       uint8  `query:"page" validate:"required,min=1"`
		Limit      uint8  `query:"limit" validate:"required,min=1,max=100"`
		Title      string `query:"title" validate:"omitempty,min=3,max=100"`
		Status     string `query:"status" validate:"omitempty,oneof='Asignado' 'Iniciado' 'En espera' 'Finalizado Éxito' 'Finalizado Error'"`
		Username   string `query:"username" validate:"omitempty,min=5,max=32"`
		DueDateMin string `query:"due_date_min" validate:"required_with=DueDateMax,omitempty,datetime=2006-01-02"`
		DueDateMax string `query:"due_date_max" validate:"required_with=DueDateMin,omitempty,datetime=2006-01-02"`
		OrderBy    string `query:"order_by" validate:"omitempty,oneof=due_date title created_at updated_at"`
	}

	UpdateRequest struct {
		Title        string `json:"title" validate:"omitempty,min=3,max=100"`
		Description  string `json:"description" validate:"omitempty,min=3,max=1000"`
		DueDate      string `json:"due_date" validate:"omitempty,datetime=2006-01-02"`
		AssignedUser string `json:"assigned_user" validate:"omitempty,uuid4"`
		Status       string `json:"status" validate:"omitempty,oneof='Iniciado' 'En espera' 'Finalizado Éxito' 'Finalizado Error'"`
	}

	CommentRequest struct {
		Comment string `json:"comment" validate:"required,min=3,max=1000"`
	}
)

type (
	StatusResponse struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	UserResponse struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}

	TaskResponse struct {
		ID           string                     `json:"id"`
		Title        string                     `json:"title"`
		Description  string                     `json:"description"`
		DueDate      string                     `json:"due_date"`
		Status       StatusResponse             `json:"status"`
		AssignedUser UserResponse               `json:"assigned_user"`
		Comments     []comments.CommentResponse `json:"comments,omitempty"`
		CreatedAt    string                     `json:"created_at"`
		UpdatedAt    string                     `json:"updated_at"`
	}

	ListResponse struct {
		Total uint           `json:"total"`
		Tasks []TaskResponse `json:"tasks"`
	}

	filterList struct {
		filters string
		args    []interface{}
	}
)

func (cr *CreateRequest) ToTaskSQL(statusID string) *TaskSQL {
	return &TaskSQL{
		Title:       cr.Title,
		Description: cr.Description,
		DueDate:     cr.DueDate,
		Status: StatusSQL{
			ID: statusID,
		},
		AssignedUser: UserSQL{
			ID: cr.AssignedUser,
		},
	}
}

func (fl *filterList) getTitle(queryParams ListRequest) *filterList {
	if queryParams.Title != "" {
		fl.args = append(fl.args, queryParams.Title)
		fl.filters += fmt.Sprintf(" AND t.title ILIKE $%d", len(fl.args))
	}

	return fl
}

func (fl *filterList) getStatus(queryParams ListRequest) *filterList {
	if queryParams.Status != "" {
		fl.args = append(fl.args, queryParams.Status)
		fl.filters += fmt.Sprintf(` AND s.current_status = $%d`, len(fl.args))
	}

	return fl
}

func (fl *filterList) getAssignedUser(queryParams ListRequest) *filterList {
	if queryParams.Username != "" {
		fl.args = append(fl.args, queryParams.Username)
		fl.filters += fmt.Sprintf(" AND u.username = $%d", len(fl.args))
	}

	return fl
}

func (fl *filterList) getDueDate(queryParams ListRequest) *filterList {
	if queryParams.DueDateMin != "" {
		fl.args = append(fl.args, queryParams.DueDateMin)
		fl.filters += fmt.Sprintf(" AND t.due_date >= $%d::date", len(fl.args))
	}
	if queryParams.DueDateMax != "" {
		fl.args = append(fl.args, queryParams.DueDateMax)
		fl.filters += fmt.Sprintf(" AND t.due_date <= $%d::date", len(fl.args))
	}

	return fl
}
