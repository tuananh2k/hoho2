package service

import (
	"context"
	"hoho-framework-v2/model"
	"hoho-framework-v2/usecase/repository"
)

type UserService interface {
	AddUser(model.User) (model.User, error)
	EditUser() (model.User, error)
	GetUser(*context.Context) (model.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{
		userRepository: r,
	}
}

func (uS *userService) AddUser(model.User) (model.User, error) {
	return uS.userRepository.Add()
}
func (uS *userService) EditUser() (model.User, error) {
	return uS.userRepository.Edit()
}
func (uS *userService) GetUser(ctx *context.Context) (model.User, error) {
	var user model.User
	user.Name = "Hoangnd"
	user.ID = "123"
	return user, nil
}
