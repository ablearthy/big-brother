package sse

import (
	"big-brother/internal/background"
	"big-brother/internal/db"
	"big-brother/internal/longpoll"
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	MIMETextEventStream = "text/event-stream"
)

type TransferLongPollMessagesService struct {
}

func (*TransferLongPollMessagesService) Transfer(c echo.Context, userId int) error {
	queries := db.New(db.GetConn())
	token, err := queries.GetTokenById(context.Background(), int32(userId))

	if err != nil || !token.Valid {
		return c.NoContent(http.StatusOK)
	}

	vt, err := queries.GetVkToken(context.Background(), token.String)

	if err != nil {
		return c.NoContent(http.StatusOK)
	}

	lpmw := background.GetLongPollManagerWrapper()

	messageCh := make(chan longpoll.Response)

	lpmw.Subscribe(&longpoll.Subscriber{
		VkUserId: longpoll.VkUserId(vt.VkUserID),
		Ch:       messageCh,
	})

	ctx := c.Request().Context()

	c.Response().Header().Set(echo.HeaderCacheControl, "no-store")
	c.Response().Header().Set(echo.HeaderContentType, MIMETextEventStream)

	select {
	case <-ctx.Done():
		return c.NoContent(http.StatusOK)
	case msg := <-messageCh:
		log.Println("New message: ", msg)

	}
	return c.NoContent(http.StatusOK)
}
