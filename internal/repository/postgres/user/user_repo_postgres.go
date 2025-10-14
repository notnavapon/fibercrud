package postgres

import (
	domainUser "clean/internal/domain/user"

	"gorm.io/gorm"
)

type userRepoPostgres struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domainUser.UserRepository {
	return &userRepoPostgres{db: db}
}

func (repo *userRepoPostgres) Create(user *domainUser.User) error {
	return repo.db.Create(user).Error
}

func (repo *userRepoPostgres) GetByEmail(email string) (*domainUser.User, error) {
	var user domainUser.User

	err := repo.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *userRepoPostgres) GetAllUser() ([]domainUser.User, error) {
	var users []domainUser.User

	if err := repo.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *userRepoPostgres) GetByID(id int) (*domainUser.User, error) {
	var user domainUser.User
	err := repo.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *userRepoPostgres) Delete(id int) error {
	err := repo.db.Unscoped().Where("id = ?", id).Delete(&domainUser.User{}).Error
	return err
}

func (repo *userRepoPostgres) Update(user *domainUser.User) (*domainUser.User, error) {
	result := repo.db.Save(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
