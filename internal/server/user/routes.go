package user

import (
	"big-brother/internal/controller/user"
	req "big-brother/internal/request/user"
	svc "big-brother/internal/service/user"
	"github.com/labstack/echo/v4"
)

func SetUserGroup(e *echo.Group) {
	cuh := user.CreateUserHandler{CreateUserService: &svc.CreateUserService{}, Validator: &req.CreateUserRequestValidator{}}
	luh := user.LoginUserHandler{
		LoginUserService: &svc.LoginUserService{},
		Validator:        &req.UserLoginRequestValidator{},
	}
	e.POST("/create", cuh.CreateUser)
	e.POST("/login", luh.Login)
}
