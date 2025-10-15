package repository

import "clean/domain/user/entity"

type UserRepository interface {
	Create(user *entity.User) error
	GetByEmail(email string) (*entity.User, error)
	GetAllUser() ([]entity.User, error)
	GetByID(id int) (*entity.User, error)
	Delete(id int) error
	Update(user *entity.User) (*entity.User, error)
}
