package controller

import (
	"hoho-framework-v2/adapters/request"
	"hoho-framework-v2/usecase/service"
	"net/http"
)

type userController struct {
	userService service.UserService
}

type UserController interface {
	AddNewUser(c *Context) error
	EditExistUser(c *Context) error
	GetExistUser(c *Context) error
}

func NewUserController(s service.UserService) UserController {
	return &userController{
		userService: s,
	}
}

func (uC *userController) AddNewUser(c *Context) error {
	rq, e := request.Make("https://syql.symper.vn/formulas/get-data").
		SetHeaders(map[string]string{"Authorization": "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjE5IiwibmFtZSI6Ik5ndXlcdTFlYzVuIFx1MDExMFx1MDBlY25oIEhvXHUwMGUwbmciLCJlbWFpbCI6ImhvYW5nbmRAc3ltcGVyLnZuIiwidGVuYW50SWQiOiIwIiwicmVzZXRQc3dkSW5mbyI6bnVsbCwidHlwZSI6ImJhIiwidXNlckRlbGVnYXRlIjp7ImlkIjoiOTcxIiwiZmlyc3ROYW1lIjoiSG9cdTAwZTBuZyIsImxhc3ROYW1lIjoiTmd1eVx1MWVjNW4gXHUwMTEwXHUwMGVjbmgiLCJ1c2VyTmFtZSI6Ik5ndXlcdTAwZWFcdTAzMDNuIFx1MDExMGlcdTAzMDBuaCBIb2FcdTAzMDBuZyIsImRpc3BsYXlOYW1lIjoiTmd1eVx1MWVjNW4gXHUwMTEwXHUwMGVjbmggSG9cdTAwZTBuZyIsImVtYWlsIjoiaG9hbmduZEBzeW1wZXIudm4iLCJ0ZW5hbnRJZCI6IjAiLCJyZXNldFBzd2RJbmZvIjpudWxsLCJ0eXBlIjoidXNlciIsImlwIjoiMTI3LjAuMC4xLCAxNC4yMjUuNDQuODAiLCJ1c2VyQWdlbnQiOiJNb3ppbGxhXC81LjAgKE1hY2ludG9zaDsgSW50ZWwgTWFjIE9TIFggMTBfMTVfNykgQXBwbGVXZWJLaXRcLzUzNy4zNiAoS0hUTUwsIGxpa2UgR2Vja28pIENocm9tZVwvMTA5LjAuMC4wIFNhZmFyaVwvNTM3LjM2Iiwicm9sZSI6Im9yZ2NoYXJ0OjExMzpmNjNkMDk4MC1hNWRkLTQyNTktODQyYS0xMDE4ODZkYjdlNTkifSwidGVuYW50Ijp7ImlkIjoiMCJ9LCJpc19jbG91ZCI6dHJ1ZSwidGVuYW50X2RvbWFpbiI6InN5bXBlci52biIsImV4cCI6MTY3Mzk3NDczNywiaWF0IjoxNjczOTcxMTM3fQ.gccGPOVHhxp1am5aGC-rmSPXU9mNuqfljtVd1qMZcTdYBot46KLbBMWDhFsXka7wvBiYvRqndBsIFATWWt8qRZJU7gcZxUDRMxXdbDUTnMZ_zKxy_3DT2uI4DLzOW4sB67Nw3wQxpu2ZP1yP7AiN55If1EDEu6Vl6cGiatk78XDUq4LLxfDE8PpwBhV6q9v6tGAXmGanxSDDOhuUf4Tk1gyCy4zt-YnSs9xpDAWAa-ExWMUP5b0CPLZsHxkKTE6LYrCl_AsvnlQvssBweOGT8VVZrhL1JMUlNflqJvYfqWPmhaQ4CSZ-7dElW0Vxef3IM_WTEAjn9y1jZm01g8uK6QlaVKA8odGMh_1OR1gbnOBMTuBEee4R2vx94-ULuJA1_KlfqWOxNOlJz58BL5EBgFN3DrDno16jV9ftZL_55ungCS5p7cuh3_X0vrchgnNYaFm0SAU5QatzCdcXnys2O2Fd7spQ_7m3YTEkdkkCCmPMKp3dmcE2B7TgH4mc6IEuQVPf4llFq0eo7k4AksdokFd7LSiUry2ElNuYk_6LIQTsbTl2oc2R8cgcNRLrUEtrc658VcyYUl5PJ6Jmxq1mHzH0l2lS-BrUnLI28izCil8nQubdqZM-_j-nKBabaJJpT3SuGs9PQ3ngeJp22BR67Cj-JpO4nKbOux0OXJGetVMnew_symper_authen_!"}).
		SetBody(map[string]interface{}{"formula": "select * from users limit 10"}).
		Post()
	// rs := rq.Data.(map[string]interface{})
	// data := rs["data"].(map[string]interface{})
	// data1 := data["data"].([](map[string]string))
	// for x,value := range data1{
	// 	s := value["create_at"]
	// }
	return c.Output(rq.Status, rq, e)
}
func (uC *userController) EditExistUser(c *Context) error {
	u, e := uC.userService.EditUser()
	return c.Output(http.StatusOK, u, e)
}

func (uC *userController) GetExistUser(c *Context) error {
	return c.Output(http.StatusOK, nil, nil)
}
