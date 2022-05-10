package index

import (
	"big-brother/internal/auth"
	"net/http"

	"github.com/labstack/echo/v4"
)

type IndexHandler struct{}

func (*IndexHandler) Index(c echo.Context) error {
	_, err := auth.GetUserId(c)
	if err != nil {
		return c.Redirect(http.StatusFound, "/login")
	}
	return c.Redirect(http.StatusFound, "/home")
}

func (*IndexHandler) Login(c echo.Context) error {
	_, err := auth.GetUserId(c)
	if err != nil {
		return c.Render(http.StatusOK, "login", struct{}{})
	}
	return c.Redirect(http.StatusFound, "/home")
}

func (*IndexHandler) Home(c echo.Context) error {
	_, err := auth.GetUserId(c)
	if err != nil {
		return c.Redirect(http.StatusFound, "/login")
	}
	return c.Render(http.StatusOK, "home", struct{}{})
}
