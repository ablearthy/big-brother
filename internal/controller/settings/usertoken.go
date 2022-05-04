package settings

import (
	"big-brother/internal/auth"
	service "big-brother/internal/interface/settings"
	req "big-brother/internal/request/settings"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserTokenSetHandler struct {
	Service   service.UserTokenSetService
	Validator *req.UserTokenSetRequestValidator
}

func (utsh *UserTokenSetHandler) SetToken(c echo.Context) error {
	request := new(req.UserTokenSetRequest)

	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := utsh.Validator.Validate(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userId, err := auth.GetUserId(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden)
	}

	if err := utsh.Service.SetToken(int32(userId), request.AccessToken); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

type UserTokenDeleteHandler struct {
	Service service.UserTokenDeleteService
}

func (utsh *UserTokenDeleteHandler) DeleteToken(c echo.Context) error {
	userId, err := auth.GetUserId(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden)
	}

	if err := utsh.Service.DeleteToken(int32(userId)); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
