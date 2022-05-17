package user

import (
	"big-brother/internal/db"
	"context"
	"errors"
)

type LastMessageUnit struct {
	ID      int         `json:"id"`
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

type LastMessageEventResponse = []LastMessageUnit

type LastMessageEventService struct{}

func (*LastMessageEventService) GetMessageEvents(userId int, lastId int) (any, error) {
	queries := db.New(db.GetConn())

	vkUserId, err := queries.GetVkUserIdByUserId(context.Background(), int32(userId))
	if err != nil || !vkUserId.Valid {
		return nil, errors.New("unable to get user info")
	}

	rows, err := queries.GetLastMessages(context.Background(), db.GetLastMessagesParams{
		VkOwnerID: vkUserId.Int32,
		ID:        int32(lastId),
		Limit:     10,
	})
	if err != nil {
		return nil, errors.New("unable to get last message events")
	}

	data := make(LastMessageEventResponse, len(rows))

	for i, row := range rows {
		data[i] = LastMessageUnit{
			ID:      int(row.ID),
			Type:    string(row.MType),
			Content: row.Message.Get(),
		}
	}

	return data, nil
}
