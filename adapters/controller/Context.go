package controller

/*
create by: Hoangnd
create at: 2023-01-01
des: Override Echo Context
*/

import (
	sAuth "hoho-framework-v2/library/auth"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Context struct {
	Context    echo.Context
	AuthObject sAuth.AuthObject
}

func (c *Context) Output(code int, i interface{}, e error) error {
	errMessage := ""
	if e != nil {
		errMessage = e.Error()
	}
	return c.Context.JSON(code, map[string]interface{}{"data": i, "error": errMessage})
}

func (c *Context) Bind(i interface{}) error {
	return c.Context.Bind(i)
}
func (c *Context) Param(i string) string {
	return c.Context.Param(i)
}
func (c *Context) FormValue(i string) string {
	return c.Context.FormValue(i)
}
func (c *Context) QueryString() string {
	return c.Context.QueryString()
}
func (c *Context) QueryParam(i string) string {
	return c.Context.QueryParam(i)
}
func (c *Context) Request() *http.Request {
	return c.Context.Request()
}
