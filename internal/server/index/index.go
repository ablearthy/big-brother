package index

import (
	"big-brother/internal/controller/index"

	"github.com/labstack/echo/v4"
)

func SetIndexGroup(e *echo.Echo) {
	ih := &index.IndexHandler{}

	e.GET("/", ih.Index)
	e.GET("/login", ih.Login)
	e.GET("/home", ih.Home)
	e.GET("/home/settings", ih.Settings)
}
