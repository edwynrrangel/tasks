package comments

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

func (r *repository) Add(comment *CommentSQL) error {
	err := r.clientDB.QueryRow(
		sqlInsert,
		comment.Task.ID,
		comment.Comment,
		comment.CreatedBy.ID,
	).
		Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		logger.Error(
			"error adding comment",
			"func", "Add - r.clientDB.QueryRow",
			"error", err.Error(),
		)
		return err
	}

	return nil
}

func (r *repository) GetByTaskID(taskID string) (comments []CommentSQL, err error) {
	query := fmt.Sprintf(sqlGet, " AND c.task_id = $1")
	if err := r.clientDB.Select(&comments, query, taskID); err != nil {
		logger.Error(
			"error getting comments by task id",
			"func", "GetByTaskID - r.clientDB.Select",
			"error", err.Error(),
		)
		return nil, err
	}

	return comments, nil
}
