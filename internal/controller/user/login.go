package user

import (
	"big-brother/internal/auth"
	service "big-brother/internal/interface/user"
	req "big-brother/internal/request/user"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type LoginUserHandler struct {
	LoginUserService service.LoginUserService
	Validator        *req.UserLoginRequestValidator
}

func (luh LoginUserHandler) Login(c echo.Context) (err error) {
	ulreq := new(req.UserLoginRequest)

	if err = c.Bind(ulreq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ulreq.Username = strings.ToLower(ulreq.Username)
	if err = luh.Validator.Validate(ulreq); err != nil {
		return err
	}
	userId, err := luh.LoginUserService.Login(ulreq)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("cannot find user in database"))
	}

	err = auth.SetCookies(c, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("failed to set cookies"))
	}

	return c.NoContent(http.StatusOK)
}
