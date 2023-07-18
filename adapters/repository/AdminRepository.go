package repository

import (
	"hoho-framework-v2/model"
	"hoho-framework-v2/usecase/repository"

	_ "github.com/lib/pq"
)

type adminRepository struct {
	db *SymperOrm
}

func NewAdminRepository(db *SymperOrm) repository.AdminRepository {
	return &adminRepository{
		db: db,
	}
}

func (u *adminRepository) Add() (model.User, error) {
	user := new(model.User)
	user.Age = "22"
	user.Name = "long admin"
	return *user, nil
}
func (u *adminRepository) Edit() (model.User, error) {
	user := new(model.User)
	user.Age = "12"
	user.Name = "hoang admin"
	return *user, nil
}
