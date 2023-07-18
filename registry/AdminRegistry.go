package registry

import (
	"hoho-framework-v2/adapters/controller"
	uR "hoho-framework-v2/adapters/repository"
	"hoho-framework-v2/usecase/repository"
	uSv "hoho-framework-v2/usecase/service"
)

func (r *registry) NewAdminController() controller.AdminController {
	return controller.NewAdminController(r.NewAdminService())
}
func (r *registry) NewAdminService() uSv.AdminService {
	return uSv.NewAdminService(r.NewAdminRepository())
}
func (r *registry) NewAdminRepository() repository.AdminRepository {
	return uR.NewAdminRepository(r.db)
}
