package longpoll

import (
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
	lp, err := sdklp.NewLongPoll(vk, 2)
	if err != nil {
		return
	}
	lp.Goroutine(true)
	go func() {
		lp.FullResponse(func(obj object.LongPollResponse) {
			lpm.publisherCh <- publisherData{
				vkUserId: mtr.VkUserId,
				obj:      obj,
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
