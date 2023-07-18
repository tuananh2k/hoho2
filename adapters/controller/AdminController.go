package controller

import (
	"hoho-framework-v2/usecase/service"
	"net/http"
)

type adminController struct {
	adminService service.AdminService
}

type AdminController interface {
	AddAdmin(c *Context) error
	EditAdmin(c *Context) error
}

func NewAdminController(s service.AdminService) AdminController {
	return &adminController{
		adminService: s,
	}
}

func (uC *adminController) AddAdmin(c *Context) error {
	u, e := uC.adminService.AddAdmin()
	return c.Output(http.StatusOK, u, e)
}
func (uC *adminController) EditAdmin(c *Context) error {
	u, e := uC.adminService.EditAdmin()
	return c.Output(http.StatusOK, u, e)
}
