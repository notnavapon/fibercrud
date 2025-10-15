package usecaseUser

import (
	dtoUser "clean/application/dto/user"
	"clean/domain/user/entity"
	domainUser "clean/domain/user/repository"
	jwtpkg "clean/infrastructure/pkg/jwt"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUsecase interface {
	CreateUser(req *dtoUser.CreateUserRequest) (*dtoUser.UserResponse, error)
	LoginUser(req *dtoUser.LoginUserRequest) (*dtoUser.LoginUserResponse, error)
	LogoutUser() error
	GetAllUser() (*dtoUser.UserListResponse, error)
	GetByID(id int) (*dtoUser.UserResponse, error)
	DeleteUser(id int) error
	UpdateUser(req *dtoUser.UpdateUserRequest) (map[string]interface{}, error)
	ChangePasswordUser(req *dtoUser.ChangePasswordRequest, email string) (*dtoUser.UserResponse, error)
}

type userUsecase struct {
	userRepo  domainUser.UserRepository
	JwtSecret string
}

func NewUserUsecase(repo domainUser.UserRepository, config string) UserUsecase {
	return &userUsecase{
		userRepo:  repo,
		JwtSecret: config,
	}
}

func (usecase *userUsecase) CreateUser(req *dtoUser.CreateUserRequest) (*dtoUser.UserResponse, error) {
	existingUser, err := usecase.userRepo.GetByEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("failed to check existing email: " + err.Error())
	}

	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, errors.New("error in hash password: " + err.Error())
	}
	fmt.Println("Password hash:", string(hashPassword))

	user := &entity.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := usecase.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &dtoUser.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (usecase *userUsecase) LoginUser(req *dtoUser.LoginUserRequest) (*dtoUser.LoginUserResponse, error) {
	checkUser, err := usecase.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("email not found")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("password don't match")
	}

	token, err := jwtpkg.GenerateToken(req.Email, []byte(usecase.JwtSecret))
	if err != nil {
		return nil, errors.New("failed to create jwt: " + err.Error())
	}

	return &dtoUser.LoginUserResponse{
		Token: token,
		User: dtoUser.UserResponse{
			ID:        checkUser.ID,
			Name:      checkUser.Name,
			Email:     checkUser.Email,
			CreatedAt: checkUser.CreatedAt.Format(time.RFC3339),
			UpdatedAt: checkUser.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (usecase *userUsecase) LogoutUser() error {
	return nil
}

func (usecase *userUsecase) GetAllUser() (*dtoUser.UserListResponse, error) {
	users, err := usecase.userRepo.GetAllUser()

	if err != nil {
		return nil, errors.New("failed to get all user: " + err.Error())
	}

	userResponse := []dtoUser.UserResponse{}

	for _, user := range users {
		userResponse = append(userResponse, dtoUser.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &dtoUser.UserListResponse{
		Users: userResponse,
	}, nil

}

func (usecase *userUsecase) GetByID(id int) (*dtoUser.UserResponse, error) {

	user, err := usecase.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &dtoUser.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (usecase *userUsecase) DeleteUser(id int) error {
	err := usecase.userRepo.Delete(id)
	return err
}

func (usecase *userUsecase) UpdateUser(req *dtoUser.UpdateUserRequest) (map[string]interface{}, error) {
	user, err := usecase.userRepo.GetByEmail(req.Email)

	if err != nil {
		return nil, errors.New("failed to update user: " + err.Error())
	}
	user.Name = req.Name
	user.Email = req.Email

	updatedUser, err := usecase.userRepo.Update(user)
	if err != nil {
		return nil, errors.New("failed to update user: " + err.Error())
	}
	res := &dtoUser.UserResponse{
		ID:        updatedUser.ID,
		Name:      updatedUser.Name,
		Email:     updatedUser.Email,
		CreatedAt: updatedUser.CreatedAt.Format(time.RFC3339),
		UpdatedAt: updatedUser.CreatedAt.Format(time.RFC3339),
	}

	token, err := jwtpkg.GenerateToken(req.Email, []byte(usecase.JwtSecret))
	if err != nil {
		return nil, errors.New("failed to create jwt: " + err.Error())
	}

	return map[string]interface{}{
		"user":  res,
		"token": token,
	}, nil
}

func (usecase *userUsecase) ChangePasswordUser(req *dtoUser.ChangePasswordRequest, email string) (*dtoUser.UserResponse, error) {
	user, err := usecase.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	//check current password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword))
	if err != nil {
		return nil, err
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashPassword)

	updatedUser, err := usecase.userRepo.Update(user)
	if err != nil {
		return nil, err
	}
	res := &dtoUser.UserResponse{
		ID:        updatedUser.ID,
		Name:      updatedUser.Name,
		Email:     updatedUser.Email,
		CreatedAt: updatedUser.CreatedAt.Format(time.RFC3339),
		UpdatedAt: updatedUser.CreatedAt.Format(time.RFC3339),
	}
	return res, nil

}
