package domainUser

type UserRepository interface {
	Create(user *User) error
	GetByEmail(email string) (*User, error)
	GetAllUser() ([]User, error)
	GetByID(id int) (*User, error)
	Delete(id int) error
	Update(user *User) (*User, error)
}
