package user

type LastMessageEventService interface {
	GetMessageEvents(userId int, lastId int) (any, error)
}
