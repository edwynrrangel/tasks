package users

import "fmt"

type (
	CreateRequest struct {
		Username       string `json:"username" validate:"required,min=5,max=32"`
		Role           string `json:"role" validate:"required,oneof=Administrator Executor Auditor"`
		PasswordLength uint8  `json:"password_length" validate:"required,min=12,max=32"`
	}

	ListRequest struct {
		Page         uint8  `query:"page" validate:"required,min=1"`
		Limit        uint8  `query:"limit" validate:"required,min=1,max=100"`
		Username     string `query:"username"`
		Role         string `query:"role" validate:"omitempty,oneof=Administrator Executor Auditor"`
		CreatedAtMin string `query:"created_at_min" validate:"required_with=CreatedAtMin,omitempty,datetime=2006-01-02"`
		CreatedAtMax string `query:"created_at_max" validate:"required_with=CreatedAtMin,omitempty,datetime=2006-01-02"`
		OrderBy      string `query:"order_by" validate:"omitempty,oneof=username role created_at"`
	}

	filterList struct {
		filters string
		args    []interface{}
	}
)

type (
	UserResponse struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		Password  string `json:"password,omitempty"`
		Role      string `json:"role"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	ListResponse struct {
		Total uint           `json:"total"`
		Users []UserResponse `json:"users"`
	}
)

func (fl *filterList) getUsername(queryParams ListRequest) *filterList {
	if queryParams.Username != "" {
		fl.args = append(fl.args, queryParams.Username)
		fl.filters += fmt.Sprintf(" AND username = $%d", len(fl.args))
	}

	return fl
}

func (fl *filterList) getRole(queryParams ListRequest) *filterList {
	if queryParams.Role != "" {
		fl.args = append(fl.args, queryParams.Role)
		fl.filters += fmt.Sprintf(` AND r.name = $%d`, len(fl.args))
	}

	return fl
}

func (fl *filterList) getCreatedAt(queryParams ListRequest) *filterList {
	if queryParams.CreatedAtMin != "" {
		fl.args = append(fl.args, queryParams.CreatedAtMin)
		fl.filters += fmt.Sprintf(" AND created_at >= $%d::date", len(fl.args))
	}
	if queryParams.CreatedAtMax != "" {
		fl.args = append(fl.args, queryParams.CreatedAtMax)
		fl.filters += fmt.Sprintf(" AND created_at <= $%d::date", len(fl.args))
	}

	return fl
}
