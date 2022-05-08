package server

import (
	"big-brother/internal/server/index"
	"big-brother/internal/server/settings"
	"big-brother/internal/server/sse"
	"big-brother/internal/server/user"

	"github.com/labstack/echo/v4"
)

func SetRoutes(e *echo.Echo) {
	index.SetIndexGroup(e)

	g := e.Group("/user")
	user.SetUserGroup(g)
	sse.SetSSEGroup(g)

	s := e.Group("/settings")
	settings.SetSettingsGroup(s)
}
