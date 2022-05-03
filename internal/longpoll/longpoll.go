package longpoll

import (
	"github.com/SevereCloud/vksdk/v2/api"
	sdklp "github.com/SevereCloud/vksdk/v2/longpoll-user"
	"github.com/SevereCloud/vksdk/v2/object"
)

type AccessToken string

type StatusRequest struct {
	at  AccessToken
	bCh chan<- bool
}

type longPollManager struct {
	clients          map[AccessToken]*sdklp.LongPoll
	acceptNewTokenCh chan AccessToken
	deleteTokenCh    chan AccessToken
	isAliveCh        chan StatusRequest
	publisherCh      chan<- publisherData
}

func initLongPollManager(publisher chan<- publisherData) *longPollManager {
	return &longPollManager{
		clients:          make(map[AccessToken]*sdklp.LongPoll),
		acceptNewTokenCh: make(chan AccessToken),
		deleteTokenCh:    make(chan AccessToken),
		isAliveCh:        make(chan StatusRequest),
		publisherCh:      publisher,
	}
}

func (lpm *longPollManager) AddNewToken(token AccessToken) {
	lpm.acceptNewTokenCh <- token
}

func (lpm *longPollManager) DeleteToken(token AccessToken) {
	lpm.deleteTokenCh <- token
}

func (lpm *longPollManager) TokenIsAlive(token AccessToken) <-chan bool {
	ch := make(chan bool)

	lpm.isAliveCh <- StatusRequest{token, ch}

	return ch
}

func (lpm *longPollManager) registerNewToken(token AccessToken) {
	if _, ok := lpm.clients[token]; ok {
		return
	}
	vk := api.NewVK(string(token))
	lp, err := sdklp.NewLongPoll(vk, 2)
	if err != nil {
		return
	}
	lp.Goroutine(true)
	go func() {
		lp.FullResponse(func(obj object.LongPollResponse) {
			lpm.publisherCh <- publisherData{
				at:  token,
				obj: obj,
			}
		})
		lp.Run()
	}()
	lpm.clients[token] = lp
}

func (lpm *longPollManager) deleteToken(token AccessToken) {
	lp, ok := lpm.clients[token]
	if !ok {
		return
	}
	lp.Shutdown()
	delete(lpm.clients, token)
}

func (lpm *longPollManager) sendAliveStatus(s *StatusRequest) {
	_, ok := lpm.clients[s.at]
	s.bCh <- ok
}

func (lpm *longPollManager) Run() {
	for {
		select {
		case newToken := <-lpm.acceptNewTokenCh:
			lpm.registerNewToken(newToken)
		case token := <-lpm.deleteTokenCh:
			lpm.deleteToken(token)
		case s := <-lpm.isAliveCh:
			lpm.sendAliveStatus(&s)
		}
	}
}
