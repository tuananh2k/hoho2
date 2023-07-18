package controller

import (
	"net/http"
)

type AppController struct {
	UserController  UserController
	AdminController AdminController
	KafkaController KafkaController
}

func (app *AppController) HelthCheck(e *Context) error {
	return e.Output(http.StatusOK, map[string]interface{}{"status": "Running"}, nil)
}
