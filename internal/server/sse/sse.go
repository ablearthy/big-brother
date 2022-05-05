package sse

import (
	handler "big-brother/internal/controller/sse"
	svc "big-brother/internal/service/sse"

	"github.com/labstack/echo/v4"
)

func SetSSEGroup(e *echo.Group) {
	sh := &handler.SSEHandler{
		Service: &svc.TransferLongPollMessagesService{},
	}
	e.GET("/sse", sh.SSE)
}
