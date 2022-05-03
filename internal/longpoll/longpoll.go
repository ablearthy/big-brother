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

type LongPollManager struct {
	clients          map[AccessToken]*sdklp.LongPoll
	acceptNewTokenCh chan AccessToken
	deleteTokenCh    chan AccessToken
	isAliveCh        chan StatusRequest
}

func Init() *LongPollManager {
	return &LongPollManager{
		clients:          make(map[AccessToken]*sdklp.LongPoll),
		acceptNewTokenCh: make(chan AccessToken),
		deleteTokenCh:    make(chan AccessToken),
		isAliveCh:        make(chan StatusRequest),
	}
}

func (lpm *LongPollManager) AddNewToken(token AccessToken) {
	lpm.acceptNewTokenCh <- token
}

func (lpm *LongPollManager) DeleteToken(token AccessToken) {
	lpm.deleteTokenCh <- token
}

func (lpm *LongPollManager) TokenIsAlive(token AccessToken) <-chan bool {
	ch := make(chan bool)

	lpm.isAliveCh <- StatusRequest{token, ch}

	return ch
}

func (lpm *LongPollManager) registerNewToken(token AccessToken) {
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
			// TODO: send messages to broker
		})
		lp.Run()
	}()
	lpm.clients[token] = lp
}

func (lpm *LongPollManager) deleteToken(token AccessToken) {
	lp, ok := lpm.clients[token]
	if !ok {
		return
	}
	lp.Shutdown()
	delete(lpm.clients, token)
}

func (lpm *LongPollManager) sendAliveStatus(s *StatusRequest) {
	_, ok := lpm.clients[s.at]
	s.bCh <- ok
}

func (lpm *LongPollManager) Run() {
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
