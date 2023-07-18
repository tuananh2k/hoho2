package repository

import (
	"hoho-framework-v2/model"
	"hoho-framework-v2/usecase/repository"

	_ "github.com/lib/pq"
)

type userRepository struct {
	db *SymperOrm
}

func NewUserRepository(db *SymperOrm) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) Add() (model.User, error) {
	user := new(model.User)
	user.Age = "12"
	user.Name = "hoang"
	// u.db.Model(user).Insert()
	return *user, nil
}
func (u *userRepository) Edit() (model.User, error) {
	user := new(model.User)
	user.Age = "13"
	// u.db.Model(user).Update()
	return *user, nil
}
func (u *userRepository) Get() (interface{}, error) {
	user := new(model.User)

	return user, nil
}
