package longpoll

import (
	"log"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	sdklp "github.com/SevereCloud/vksdk/v2/longpoll-user"
	"github.com/SevereCloud/vksdk/v2/object"
)

type AccessToken string
type VkUserId int32

type ManageTokenRequest struct {
	AccessToken AccessToken
	VkUserId    VkUserId
}

type StatusRequest struct {
	vkUserId VkUserId
	bCh      chan<- bool
}

type LongPoll struct {
	lp      *sdklp.LongPoll
	counter uint32
}

type longPollManager struct {
	clients          map[VkUserId]*LongPoll
	acceptNewTokenCh chan ManageTokenRequest
	deleteTokenCh    chan ManageTokenRequest
	isAliveCh        chan StatusRequest
	publisherCh      chan<- publisherData
}

func initLongPollManager(publisher chan<- publisherData) *longPollManager {
	return &longPollManager{
		clients:          make(map[VkUserId]*LongPoll),
		acceptNewTokenCh: make(chan ManageTokenRequest),
		deleteTokenCh:    make(chan ManageTokenRequest),
		isAliveCh:        make(chan StatusRequest),
		publisherCh:      publisher,
	}
}

func (lpm *longPollManager) AddNewToken(mtr ManageTokenRequest) {
	lpm.acceptNewTokenCh <- mtr
}

func (lpm *longPollManager) DeleteToken(mtr ManageTokenRequest) {
	lpm.deleteTokenCh <- mtr
}

func (lpm *longPollManager) TokenIsAlive(vkUserId VkUserId) <-chan bool {
	ch := make(chan bool)

	lpm.isAliveCh <- StatusRequest{vkUserId, ch}

	return ch
}

func (lpm *longPollManager) registerNewToken(mtr ManageTokenRequest) {
	if v, ok := lpm.clients[mtr.VkUserId]; ok {
		v.counter += 1
		return
	}
	vk := api.NewVK(string(mtr.AccessToken))
	lp, err := sdklp.NewLongPoll(vk, 234)
	if err != nil {
		return
	}
	lp.Goroutine(true)
	go func() {
		lp.FullResponse(func(obj object.LongPollResponse) {
			for _, ev := range obj.Updates {
				eventType := int(ev[0].(float64))
				switch eventType {
				case 2:
					messageId := int(ev[1].(float64))
					flags := int(ev[2].(float64))
					if (flags & FLAG_DELETE_FOR_ALL) == 0 {
						continue
					}
					lpm.publisherCh <- formPublisherData(mtr.VkUserId, EventDeleteMessage{
						MessageId: messageId,
					})
					log.Println("Flags changed:", messageId, flags)
				case 4:
					messageId := int(ev[1].(float64))
					fullMessage := new(MessagesGetByIDExtendedResponse)
					err := lp.VK.RequestUnmarshal("messages.getById", fullMessage, api.Params{
						"message_ids": []int{messageId},
						"extended":    1,
					})
					if err != nil || fullMessage.Count != 1 {
						log.Println("New message error!")
						continue
					}
					message := Message{
						Message:  fullMessage.Items[0],
						Profiles: fullMessage.Profiles,
						Groups:   fullMessage.Groups,
					}
					lpm.publisherCh <- formPublisherData(mtr.VkUserId, EventNewMessage{
						MessageId: messageId,
						Message:   message,
					})
					log.Println("New message:", messageId)
				case 5:
					messageId := int(ev[1].(float64))
					fullMessage := new(MessagesGetByIDExtendedResponse)
					err := lp.VK.RequestUnmarshal("messages.getById", fullMessage, api.Params{
						"message_ids": []int{messageId},
						"extended":    1,
					})
					if err != nil || fullMessage.Count != 1 {
						log.Println("New message error!")
						continue
					}
					message := Message{
						Message:  fullMessage.Items[0],
						Profiles: fullMessage.Profiles,
						Groups:   fullMessage.Groups,
					}
					lpm.publisherCh <- formPublisherData(mtr.VkUserId, EventEditMessage{
						MessageId: messageId,
						Message:   message,
					})
					log.Println("Edit message:", messageId)
				case 8:
					userId := -int(ev[1].(float64))
					platformType := int(ev[2].(float64))
					lpm.publisherCh <- formPublisherData(mtr.VkUserId, EventFriendOnline{
						UserId:    VkUserId(userId),
						Platform:  int32(platformType),
						Timestamp: time.Now(),
					})
					log.Println("User online:", userId, platformType)
				case 9:
					userId := -int(ev[1].(float64))
					kickedByTimeout := int(ev[2].(float64)) != 0
					lpm.publisherCh <- formPublisherData(mtr.VkUserId, EventFriendOffline{
						UserId:          VkUserId(userId),
						KickedByTimeout: kickedByTimeout,
						Timestamp:       time.Now(),
					})
					log.Println("User offline:", userId, kickedByTimeout)
				}
				log.Println("Event:", ev)
			}
		})
		lp.Run()
		lpm.DeleteToken(mtr)
	}()
	lpm.clients[mtr.VkUserId] = &LongPoll{
		lp:      lp,
		counter: 1,
	}
}

func (lpm *longPollManager) deleteToken(mtr ManageTokenRequest) {
	lp, ok := lpm.clients[mtr.VkUserId]
	if !ok {
		return
	}
	lp.counter -= 1
	if lp.counter == 0 {
		lp.lp.Shutdown()
		delete(lpm.clients, mtr.VkUserId)
	}
}

func (lpm *longPollManager) sendAliveStatus(s *StatusRequest) {
	_, ok := lpm.clients[s.vkUserId]
	s.bCh <- ok
}

func (lpm *longPollManager) Run() {
	for {
		select {
		case req := <-lpm.acceptNewTokenCh:
			lpm.registerNewToken(req)
		case req := <-lpm.deleteTokenCh:
			lpm.deleteToken(req)
		case s := <-lpm.isAliveCh:
			lpm.sendAliveStatus(&s)
		}
	}
}
