package comments

type (
	UserSQL struct {
		ID       string `db:"id"`
		Username string `db:"username"`
	}

	TaskSQL struct {
		ID    string `db:"id"`
		Title string `db:"title"`
	}

	CommentSQL struct {
		ID        string   `db:"id"`
		Task      *TaskSQL `db:"task_id"`
		Comment   string   `db:"comment"`
		CreatedBy UserSQL  `db:"created_by"`
		CreatedAt string   `db:"created_at"`
	}
)

func (cs *CommentSQL) ToResponse() *CommentResponse {
	task := new(TaskResponse)
	if cs.Task != nil {
		task = &TaskResponse{
			ID:    cs.Task.ID,
			Title: cs.Task.Title,
		}
	}

	comment := &CommentResponse{
		ID:      cs.ID,
		Comment: cs.Comment,
		CreatedBy: UserResponse{
			ID:       cs.CreatedBy.ID,
			Username: cs.CreatedBy.Username,
		},
		CreatedAt: cs.CreatedAt,
		Task:      task,
	}

	return comment
}
