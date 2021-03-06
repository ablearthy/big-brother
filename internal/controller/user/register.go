package user

import (
	service "big-brother/internal/interface/user"
	req "big-brother/internal/request/user"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type CreateUserHandler struct {
	CreateUserService service.CreateUserService
	Validator         *req.CreateUserRequestValidator
}

func (cuh CreateUserHandler) CreateUser(c echo.Context) (err error) {
	cureq := new(req.CreateUserRequest)

	if err = c.Bind(cureq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	cureq.Username = strings.ToLower(cureq.Username)
	if err = cuh.Validator.Validate(cureq); err != nil {
		return err
	}
	if err = cuh.CreateUserService.RegisterUser(cureq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
