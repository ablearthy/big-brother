package longpoll

import "time"

const FLAG_DELETE_FOR_ALL = 131072

type EventFriendOnline struct {
	UserId    VkUserId
	Platform  int32
	Timestamp time.Time
}

type EventFriendOffline struct {
	UserId          VkUserId
	KickedByTimeout bool
	Timestamp       time.Time
}

type EventNewMessage struct {
	MessageId int
	Message   Message
}

type EventEditMessage struct {
	MessageId int
	Message   Message
}

type EventDeleteMessage struct {
	MessageId int
}

type EventType interface {
	EventFriendOffline | EventFriendOnline | EventNewMessage | EventEditMessage | EventDeleteMessage
}

type EventWrapper[T EventType] struct {
	UserId VkUserId
	Event  T
}

type MessagesGetByIDExtendedResponse struct {
	Count    int           `json:"count"`
	Items    []interface{} `json:"items"`
	Profiles []interface{} `json:"profiles,omitempty"`
	Groups   []interface{} `json:"groups,omitempty"`
}

type Message struct {
	Message  interface{}   `json:"message"`
	Profiles []interface{} `json:"profiles,omitempty"`
	Groups   []interface{} `json:"groups,omitempty"`
}

func formPublisherData[T EventType](userId VkUserId, ev T) publisherData {
	return publisherData{
		vkUserId: userId,
		obj: EventWrapper[T]{
			UserId: userId,
			Event:  ev,
		},
	}
}
