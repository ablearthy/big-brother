package dbsave

import (
	"big-brother/internal/db"
	"big-brother/internal/longpoll"
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/jackc/pgtype"
)

type DbLongPollSaver struct {
	ch chan longpoll.Response
}

func New() *DbLongPollSaver {
	return &DbLongPollSaver{
		ch: make(chan longpoll.Response),
	}
}

func (dblps *DbLongPollSaver) GetChannel() chan longpoll.Response {
	return dblps.ch
}

func (dblps *DbLongPollSaver) Run() {
	for {
		data := <-dblps.ch
		switch v := data.(type) {
		case longpoll.EventWrapper[longpoll.EventDeleteMessage]:
			log.Println("DBSAVE: [delete]:", v.UserId, v.Event.MessageId)
			err := saveEventDeleteMessage(&v)
			if err != nil {
				log.Println("DBSAVE: [delete]: an error occured while saving event:", err)
			}
		case longpoll.EventWrapper[longpoll.EventNewMessage]:
			log.Println("DBSAVE: [new]:", v.UserId, v.Event.MessageId)
			err := saveEventNewMessage(&v)
			if err != nil {
				log.Println("DBSAVE: [new]: an error occured while saving event:", err)
			}
		case longpoll.EventWrapper[longpoll.EventEditMessage]:
			log.Println("DBSAVE: [edit]:", v.UserId, v.Event.MessageId)
			err := saveEventEditMessage(&v)
			if err != nil {
				log.Println("DBSAVE: [edit]: an error occured while saving event:", err)
			}
		case longpoll.EventWrapper[longpoll.EventFriendOffline]:
			log.Println("DBSAVE: [offline]:", v.UserId, v.Event.UserId, v.Event.Timestamp)
			err := saveEventOffline(&v)
			if err != nil {
				log.Println("DBSAVE: [offline]: an error occured while saving event:", err)
			}
		case longpoll.EventWrapper[longpoll.EventFriendOnline]:
			log.Println("DBSAVE: [online]:", v.UserId, v.Event.UserId, v.Event.Platform, v.Event.Timestamp)
			err := saveEventOnline(&v)
			if err != nil {
				log.Println("DBSAVE: [online]: an error occured while saving event:", err)
			}
		}
		log.Println("DBSAVE:", data)
	}
}

func saveVKMessage(q *db.Queries, vkOwnerId int32, messageId int32, msg longpoll.Message) (int32, error) {
	m := pgtype.JSONB{}
	err := m.Set(msg)
	if err != nil {
		return 0, err
	}
	id, err := q.SaveVkMessage(context.Background(), db.SaveVkMessageParams{
		VkOwnerID: vkOwnerId,
		MessageID: messageId,
		Message:   m,
	})
	return id, err
}

func saveEventNewMessage(ev *longpoll.EventWrapper[longpoll.EventNewMessage]) error {
	tx, err := db.GetConn().Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	q := db.New(tx)
	internal_id, err := saveVKMessage(q, int32(ev.UserId), int32(ev.Event.MessageId), ev.Event.Message)
	if err != nil {
		return err
	}
	err = q.SaveMessageEvent(context.Background(), db.SaveMessageEventParams{
		InternalMessageID: internal_id,
		MType:             db.VkMessageEventTypeNew,
		CreatedAt:         time.Now(),
	})
	if err != nil {
		return err
	}
	return tx.Commit(context.Background())
}

func saveEventEditMessage(ev *longpoll.EventWrapper[longpoll.EventEditMessage]) error {
	tx, err := db.GetConn().Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	q := db.New(tx)
	internal_id, err := saveVKMessage(q, int32(ev.UserId), int32(ev.Event.MessageId), ev.Event.Message)
	if err != nil {
		return err
	}
	err = q.SaveMessageEvent(context.Background(), db.SaveMessageEventParams{
		InternalMessageID: internal_id,
		MType:             db.VkMessageEventTypeEdit,
		CreatedAt:         time.Now(),
	})
	if err != nil {
		return err
	}
	return tx.Commit(context.Background())
}

func saveEventDeleteMessage(ev *longpoll.EventWrapper[longpoll.EventDeleteMessage]) error {
	q := db.New(db.GetConn())
	lastMsgId, err := q.GetLastSavedVKMessage(context.Background(), db.GetLastSavedVKMessageParams{
		VkOwnerID: int32(ev.UserId),
		MessageID: int32(ev.Event.MessageId),
	})
	if err != nil {
		return err
	}
	if lastMsgId == nil {
		tx, err := db.GetConn().Begin(context.Background())
		if err != nil {
			return err
		}
		defer tx.Rollback(context.Background())

		queries := db.New(tx)

		m := pgtype.JSONB{}
		err = m.Set(struct {
			Error string `json:"error"`
		}{
			Error: "no_content",
		})
		if err != nil {
			return err
		}
		internalId, err := queries.SaveVkMessage(context.Background(), db.SaveVkMessageParams{
			VkOwnerID: int32(ev.UserId),
			MessageID: int32(ev.Event.MessageId),
			Message:   m,
		})
		if err != nil {
			return err
		}

		err = queries.SaveMessageEvent(context.Background(), db.SaveMessageEventParams{
			InternalMessageID: internalId,
			MType:             db.VkMessageEventTypeDelete,
			CreatedAt:         time.Now(),
		})
		if err != nil {
			return err
		}

		return tx.Commit(context.Background())
	}
	msgId := lastMsgId.(int32)

	err = q.SaveMessageEvent(context.Background(), db.SaveMessageEventParams{
		InternalMessageID: msgId,
		MType:             db.VkMessageEventTypeDelete,
		CreatedAt:         time.Now(),
	})
	return err
}

func saveEventOffline(ev *longpoll.EventWrapper[longpoll.EventFriendOffline]) error {
	q := db.New(db.GetConn())
	err := q.SaveActivityEvent(context.Background(), db.SaveActivityEventParams{
		VkOwnerID: int32(ev.UserId),
		TargetID:  int32(ev.Event.UserId),
		Activity:  db.VkActivityOffline,
		Platform:  db.VkPlatformWeb,
		KickedByTimeout: sql.NullBool{
			Bool:  ev.Event.KickedByTimeout,
			Valid: true,
		},
		CreatedAt: ev.Event.Timestamp,
	})
	return err
}

func transformPlatform(platform int32) db.VkPlatform {
	switch platform {
	case 1:
		return db.VkPlatformMobile
	case 2:
		return db.VkPlatformIphone
	case 3:
		return db.VkPlatformIpad
	case 4:
		return db.VkPlatformAndroid
	case 5:
		return db.VkPlatformWphone
	case 6:
		return db.VkPlatformWindows
	default:
		return db.VkPlatformWeb
	}
}

func saveEventOnline(ev *longpoll.EventWrapper[longpoll.EventFriendOnline]) error {
	q := db.New(db.GetConn())
	err := q.SaveActivityEvent(context.Background(), db.SaveActivityEventParams{
		VkOwnerID: int32(ev.UserId),
		TargetID:  int32(ev.Event.UserId),
		Activity:  db.VkActivityOnline,
		Platform:  transformPlatform(ev.Event.Platform),
		KickedByTimeout: sql.NullBool{
			Bool:  false,
			Valid: false,
		},
		CreatedAt: ev.Event.Timestamp,
	})
	return err
}
