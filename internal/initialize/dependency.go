package initialize

import (
	handlerUser "clean/internal/delivery/user"
	postgres "clean/internal/repository/postgres/user"
	usecaseUser "clean/internal/usecase/user"

	"gorm.io/gorm"
)

type Dependencies struct {
	UserHandler *handlerUser.UserHandler
}

func NewDependencies(db *gorm.DB, config string) *Dependencies {
	// Repository
	userRepo := postgres.NewUserRepository(db)

	// Use Case
	userUsecase := usecaseUser.NewUserUsecase(userRepo, config)

	// Handler
	userHandler := handlerUser.NewUserHandler(userUsecase)

	return &Dependencies{
		UserHandler: userHandler,
	}
}
