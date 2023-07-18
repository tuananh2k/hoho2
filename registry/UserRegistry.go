package registry

import (
	"hoho-framework-v2/adapters/controller"
	uR "hoho-framework-v2/adapters/repository"
	"hoho-framework-v2/usecase/repository"
	uSv "hoho-framework-v2/usecase/service"
)

func (r *registry) NewUserController() controller.UserController {
	return controller.NewUserController(r.NewUserService())
}
func (r *registry) NewUserService() uSv.UserService {
	return uSv.NewUserService(r.NewUserRepository())
}
func (r *registry) NewUserRepository() repository.UserRepository {
	return uR.NewUserRepository(r.db)
}
