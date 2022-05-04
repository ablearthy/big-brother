package settings

import (
	"big-brother/internal/controller/settings"

	validator "big-brother/internal/request/settings"
	service "big-brother/internal/service/settings"

	"github.com/labstack/echo/v4"
)

func SetUserTokenGroup(e *echo.Group) {
	sh := settings.UserTokenSetHandler{
		Service:   &service.UserTokenSetService{},
		Validator: &validator.UserTokenSetRequestValidator{},
	}
	e.POST("/set", sh.SetToken)
	dh := settings.UserTokenDeleteHandler{
		Service: &service.UserTokenDeleteService{},
	}
	e.POST("/delete", dh.DeleteToken)
}
