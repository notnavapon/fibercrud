package initialize

import (
	usecaseUser "clean/application/usecase/user"
	handlerUser "clean/infrastructure/http/delivery/user"
	postgres "clean/infrastructure/repository/postgres/user"

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
