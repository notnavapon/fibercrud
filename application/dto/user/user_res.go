package dtoUser

type UserResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type LoginUserResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type UserListResponse struct {
	Users []UserResponse `json:"users"`
}
