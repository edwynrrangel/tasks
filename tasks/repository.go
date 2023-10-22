package tasks

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

func (r *repository) Create(task *TaskSQL) error {
	err := r.clientDB.
		QueryRow(sqlInsert, task.Title, task.Description, task.DueDate, task.Status.ID, task.AssignedUser.ID).
		Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		logger.Error(
			"error creating task",
			"func", "Create - r.clientDB.Exec",
			"error", err.Error(),
		)
		return err
	}

	return nil
}

func (r *repository) List(queryParams ListRequest) (tasks []TaskSQL, total uint, err error) {
	fl := new(filterList)
	fl.getDueDate(queryParams).
		getStatus(queryParams).
		getTitle(queryParams).
		getAssignedUser(queryParams)

	query := fmt.Sprintf(sqlCountTasks, fl.filters)
	if err = r.clientDB.Get(&total, query, fl.args...); err != nil {
		logger.Error(
			"error getting total tasks",
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

	query = fmt.Sprintf(sqlGetTask, fl.filters)
	if err = r.clientDB.Select(&tasks, query, fl.args...); err != nil {
		logger.Error(
			"error getting tasks",
			"func", "List - r.clientDB.Select",
			"error", err.Error(),
		)
		return nil, total, err
	}

	return
}

func (r *repository) GetByID(id string) (task *TaskSQL, err error) {
	task = new(TaskSQL)
	query := fmt.Sprintf(sqlGetTask, " AND t.id = $1")
	if err = r.clientDB.Get(task, query, id); err != nil {
		logger.Error(
			"error getting task by id",
			"func", "GetByID - r.clientDB.Get",
			"error", err.Error(),
		)
		return nil, err
	}

	return
}

func (r *repository) Update(id string, task *TaskSQL) error {
	query := `
		UPDATE task.tasks
		SET 
			title = $1 
			,description = $2
			,due_date = $3 
			,assigned_user = $4
			,status = $5
			,updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at
	`

	err := r.clientDB.QueryRowx(
		query,
		task.Title,
		task.Description,
		task.DueDate,
		task.AssignedUser.ID,
		task.Status.ID,
		id,
	).Scan(&task.UpdatedAt)
	if err != nil {
		logger.Error(
			"error updating task",
			"func", "Update - r.clientDB.Exec",
			"error", err.Error(),
		)
		return err
	}

	return nil
}

func (r *repository) Delete(id string) error {
	query := `
		DELETE FROM task.tasks
		WHERE id = $1
	`

	result, err := r.clientDB.Exec(query, id)
	if err != nil {
		logger.Error(
			"error deleting task",
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
