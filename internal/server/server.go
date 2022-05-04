package server

import (
	"big-brother/internal/server/settings"
	"big-brother/internal/server/user"

	"github.com/labstack/echo/v4"
)

func SetRoutes(e *echo.Echo) {
	g := e.Group("/user")
	user.SetUserGroup(g)

	s := e.Group("/settings")
	settings.SetSettingsGroup(s)
}
