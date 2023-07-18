package registry

import (
	"hoho-framework-v2/adapters/controller"
	"hoho-framework-v2/adapters/repository"
	"hoho-framework-v2/infrastructure/cache"
)

type registry struct {
	db  *repository.SymperOrm
	rdb cache.RedisClient
}

type Registry interface {
	NewAppController() controller.AppController
}

func NewRegistry(db *repository.SymperOrm, rdb cache.RedisClient) Registry {
	return &registry{
		db:  db,
		rdb: rdb,
	}
}

func (r *registry) NewAppController() controller.AppController {
	return controller.AppController{
		UserController:  r.NewUserController(),
		AdminController: r.NewAdminController(),
		KafkaController: r.NewKafkaController(),
	}
}
