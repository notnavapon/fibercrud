package postgres

import (
	"clean/domain/user/entity"
	"clean/domain/user/repository"

	"gorm.io/gorm"
)

type userRepoPostgres struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepoPostgres{db: db}
}

func (repo *userRepoPostgres) Create(user *entity.User) error {
	return repo.db.Create(user).Error
}

func (repo *userRepoPostgres) GetByEmail(email string) (*entity.User, error) {
	var user entity.User

	err := repo.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *userRepoPostgres) GetAllUser() ([]entity.User, error) {
	var users []entity.User

	if err := repo.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *userRepoPostgres) GetByID(id int) (*entity.User, error) {
	var user entity.User
	err := repo.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *userRepoPostgres) Delete(id int) error {
	err := repo.db.Unscoped().Where("id = ?", id).Delete(&entity.User{}).Error
	return err
}

func (repo *userRepoPostgres) Update(user *entity.User) (*entity.User, error) {
	result := repo.db.Save(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
