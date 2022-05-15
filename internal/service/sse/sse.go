package sse

import (
	"big-brother/internal/background"
	"big-brother/internal/db"
	"big-brother/internal/longpoll"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	MIMETextEventStream = "text/event-stream"
)

type TransferLongPollMessagesService struct {
}

func sendSSE(w interface {
	io.Writer
	http.Flusher
}, eventName string, data any) error {
	_, err := fmt.Fprintf(w, "event: %s\ndata: ", eventName)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(data)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "\n")
	if err != nil {
		return err
	}
	w.Flush()
	return nil
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
	sub := &longpoll.Subscriber{
		VkUserId: longpoll.VkUserId(vt.VkUserID),
		Ch:       messageCh,
	}
	lpmw.Subscribe(sub)
	defer lpmw.Unsubscribe(sub)

	ctx := c.Request().Context()

	c.Response().Header().Set(echo.HeaderCacheControl, "no-store")
	c.Response().Header().Set(echo.HeaderContentType, MIMETextEventStream)

	pingTicker := time.NewTicker(15 * time.Second)
	for {
		select {
		case <-ctx.Done():
			pingTicker.Stop()
			return c.NoContent(http.StatusOK)
		case msg := <-messageCh:
			switch v := msg.(type) {
			case longpoll.EventWrapper[longpoll.EventDeleteMessage]:
				queries := db.New(db.GetConn())
				internalId, err := queries.GetLastSavedVKMessage(context.Background(), db.GetLastSavedVKMessageParams{
					VkOwnerID: int32(v.UserId),
					MessageID: int32(v.Event.MessageId),
				})
				if err != nil {
					sendSSE(c.Response(), "delete_message", "error")
					continue
				}
				if internalId == nil {
					sendSSE(c.Response(), "delete_message", "error_no_content")
					continue
				}
				msg, err := queries.GetMessageById(context.Background(), internalId.(int32))
				if err != nil {
					sendSSE(c.Response(), "delete_message", "error")
					continue
				}
				var m any
				err = json.Unmarshal(msg.Message.Bytes, &m)
				if err != nil {
					sendSSE(c.Response(), "delete_message", "error")
					continue
				}
				sendSSE(c.Response(), "delete_message", struct {
					FromID  int         `json:"from_id"`
					Content interface{} `json:"content"`
				}{
					FromID:  int(msg.MessageID),
					Content: m,
				})

			case longpoll.EventWrapper[longpoll.EventEditMessage]:
				sendSSE(c.Response(), "edit_message", v.Event)
			case longpoll.EventWrapper[longpoll.EventNewMessage]:
				sendSSE(c.Response(), "new_message", v.Event)
			}
		case <-pingTicker.C:
			sendSSE(c.Response(), "ping", struct{}{})
		}
	}
}
