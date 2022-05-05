package sse

import "github.com/labstack/echo/v4"

type TransferLongPollMessagesService interface {
	Transfer(c echo.Context, userId int) error
}
