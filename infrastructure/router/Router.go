package router

/*
create by: Hoangnd
create at: 2023-01-01
des: Xử lý router & authen
*/

import (
	"encoding/json"
	"fmt"
	"hoho-framework-v2/adapters/controller"
	"hoho-framework-v2/library"
	aAuth "hoho-framework-v2/library/auth"
	"os"
	"reflect"
	"runtime"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(e *echo.Echo, c controller.AppController) *echo.Echo {
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	var authObject aAuth.AuthObject
	config := getMiddleWareConfig(&authObject)
	e.Use(middleware.JWTWithConfig(config))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ct := &controller.Context{
				Context: c,
			}
			return next(ct.Context)
		}
	})
	e.Static("/", "static/index.html")
	e.GET("/stop-subscribe-kafka/:processId", func(context echo.Context) error { return forward(context, authObject, c.KafkaController.StopProcess) })
	e.POST("/subscribe-kafka", func(context echo.Context) error { return forward(context, authObject, c.KafkaController.Subscribe) })
	e.GET("/helthcheck", func(context echo.Context) error { return forward(context, authObject, c.HelthCheck) })
	e.GET("/user", func(context echo.Context) error { return forward(context, authObject, c.UserController.GetExistUser) })
	e.POST("/user", func(context echo.Context) error { return forward(context, authObject, c.UserController.AddNewUser) })
	e.PUT("/user", func(context echo.Context) error { return forward(context, authObject, c.UserController.EditExistUser) })
	e.PUT("/admin", func(context echo.Context) error { return forward(context, authObject, c.AdminController.EditAdmin) })
	e.POST("/admin", func(context echo.Context) error { return forward(context, authObject, c.AdminController.AddAdmin) })

	return e
}
func forward(context echo.Context, authObject aAuth.AuthObject, f func(*controller.Context) error) error {
	ct := &controller.Context{}
	ct.Context = context
	ct.AuthObject = authObject
	before := library.GetTimeMiliseconds("")
	fun := f(ct)
	errorFunc := ""
	if fun != nil {
		errorFunc = fun.Error()
	}
	after := library.GetTimeMiliseconds("")
	functionName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	userId := authObject.GetUserId()
	userRole := authObject.GetUserRole()
	tenantId := authObject.GetUserTenantId()
	tenantIdStr := strconv.Itoa(tenantId)
	dataLog := map[string]interface{}{
		"parameters":  "",
		"method":      context.Request().Method,
		"action":      functionName,
		"requestTime": after - before,
		"uri":         context.Request().URL.Path,
		"host":        context.Request().Host,
		"queryString": context.Request().URL.RawQuery,
		"userAgent":   context.Request().UserAgent(),
		"clientIp":    context.RealIP(),
		"serverIp":    "",
		"timeStamp":   library.GetCurrentTimeStamp(),
		"date":        library.GetCurrentDate(),
		"error":       errorFunc,
		"statusCode":  0,
		"output":      "",
		"userId":      userId,
		"userRole":    userRole,
		"tenantId":    tenantIdStr,
	}
	out, _ := json.Marshal(dataLog)
	file, err := os.OpenFile("log/request-"+library.GetCurrentDate()+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println("Open File Errors", err)
	} else {
		defer file.Close()
		if _, err := file.Write([]byte(string(out) + "\n")); err != nil {
			fmt.Println("Write File Errors", err)
		}
	}

	return fun
}
