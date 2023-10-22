package tasks

import (
	"time"
)

type (
	StatusSQL struct {
		ID   string `db:"id"`
		Name string `db:"name"`
	}

	UserSQL struct {
		ID       string `db:"id"`
		Username string `db:"username"`
	}

	TaskSQL struct {
		ID           string    `db:"id"`
		Title        string    `db:"title"`
		Description  string    `db:"description"`
		DueDate      string    `db:"due_date"`
		Status       StatusSQL `db:"status"`
		AssignedUser UserSQL   `db:"assigned_user"`
		CreatedAt    string    `db:"created_at"`
		UpdatedAt    string    `db:"updated_at"`
	}
)

func (ts *TaskSQL) ToResponse() *TaskResponse {
	return &TaskResponse{
		ID:          ts.ID,
		Title:       ts.Title,
		Description: ts.Description,
		DueDate:     ts.DueDate,
		CreatedAt:   ts.CreatedAt,
		UpdatedAt:   ts.UpdatedAt,
		Status: StatusResponse{
			ID:   ts.Status.ID,
			Name: ts.Status.Name,
		},
		AssignedUser: UserResponse{
			ID:       ts.AssignedUser.ID,
			Username: ts.AssignedUser.Username,
		},
	}
}

func (ts *TaskSQL) updateTitle(body *UpdateRequest) *TaskSQL {
	if body.Title != "" && body.Title != ts.Title {
		ts.Title = body.Title
	}

	return ts
}

func (ts *TaskSQL) updateDescription(body *UpdateRequest) *TaskSQL {
	if body.Description != "" && body.Description != ts.Description {
		ts.Description = body.Description
	}

	return ts
}

func (ts *TaskSQL) updateDueDate(body *UpdateRequest) *TaskSQL {
	if body.DueDate != "" && body.DueDate != ts.DueDate {
		ts.DueDate = body.DueDate
	}

	return ts
}

func (ts *TaskSQL) updateAssignedUser(body *UpdateRequest) *TaskSQL {
	if body.AssignedUser != "" && body.AssignedUser != ts.AssignedUser.ID {
		ts.AssignedUser.ID = body.AssignedUser
	}

	return ts
}

func (ts *TaskSQL) updateStatus(states StatusSQL) *TaskSQL {
	ts.Status.ID = states.ID
	ts.Status.Name = states.Name

	return ts
}

func (ts *TaskSQL) IsExpired() bool {
	dueDate, err := time.Parse(time.RFC3339, ts.DueDate)
	if err != nil {
		return false
	}

	currentDate := time.Now().UTC().Truncate(24 * time.Hour)
	return dueDate.Before(currentDate)
}
