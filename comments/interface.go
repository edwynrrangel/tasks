package comments

type Repository interface {
	Add(comment *CommentSQL) error
	GetByTaskID(taskID string) ([]CommentSQL, error)
}
