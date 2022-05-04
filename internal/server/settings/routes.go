package settings

import "github.com/labstack/echo/v4"

func SetSettingsGroup(e *echo.Group) {
	g := e.Group("/usertoken")
	SetUserTokenGroup(g)
}
