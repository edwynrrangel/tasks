package comments

type (
	TaskResponse struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}

	UserResponse struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}

	CommentResponse struct {
		ID        string        `json:"id"`
		Task      *TaskResponse `json:"task_id,omitempty"`
		Comment   string        `json:"comment"`
		CreatedBy UserResponse  `json:"created_by"`
		CreatedAt string        `json:"created_at"`
	}
)
