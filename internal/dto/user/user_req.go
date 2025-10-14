package dtoUser

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentpassword"`
	NewPassword     string `json:"newpassword"`
}

type UpdateUserRequest struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	CurrentEmail string `json:"currentemail"`
}
