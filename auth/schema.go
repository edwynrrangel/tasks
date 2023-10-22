package auth

type (
	LoginRequest struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password"`
	}

	ChangePasswordRequest struct {
		CurrentPassword    string `json:"current_password" validate:"required"`
		NewPassword        string `json:"new_password" validate:"required,min=8"`
		ConfirmNewPassword string `json:"confirm_new_password" validate:"required,eqfield=NewPassword"`
	}
)

type (
	LoginResponse struct {
		AccessToken string `json:"access_token,omitempty"`
	}
)
