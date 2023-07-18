package repository

import "hoho-framework-v2/model"

type AdminRepository interface {
	Add() (model.User, error)
	Edit() (model.User, error)
}
