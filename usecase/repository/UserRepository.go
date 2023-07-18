package repository

import "hoho-framework-v2/model"

type UserRepository interface {
	Add() (model.User, error)
	Edit() (model.User, error)
	Get() (interface{}, error)
}
