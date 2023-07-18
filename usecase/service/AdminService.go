package service

import (
	"hoho-framework-v2/model"
	"hoho-framework-v2/usecase/repository"
)

type AdminService interface {
	AddAdmin() (model.User, error)
	EditAdmin() (model.User, error)
}

type adminService struct {
	adminRepository repository.AdminRepository
}

func NewAdminService(r repository.AdminRepository) AdminService {
	return &adminService{
		adminRepository: r,
	}
}

func (uS *adminService) AddAdmin() (model.User, error) {
	return uS.adminRepository.Add()
}
func (uS *adminService) EditAdmin() (model.User, error) {
	return uS.adminRepository.Add()
}
