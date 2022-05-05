package sse

import (
	"big-brother/internal/auth"
	svc "big-brother/internal/interface/sse"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SSEHandler struct {
	Service svc.TransferLongPollMessagesService
}

func (sh *SSEHandler) SSE(c echo.Context) error {
	userId, err := auth.GetUserId(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	return sh.Service.Transfer(c, userId)
}
