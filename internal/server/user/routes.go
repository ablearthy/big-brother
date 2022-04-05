package user

import (
	"big-brother/internal/controller/user"
	req "big-brother/internal/request/user"
	svc "big-brother/internal/service/user"
	"github.com/labstack/echo/v4"
)

func SetUserGroup(e *echo.Group) {
	cuh := user.CreateUserHandler{CreateUserService: &svc.CreateUserService{}, Validator: &req.CreateUserRequestValidator{}}
	e.POST("/create", cuh.CreateUser)
}
