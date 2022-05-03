package longpoll

type LongPollManagerWrapper struct {
	lpm *longPollManager
	lpb *longPollBroker
}

func Init() *LongPollManagerWrapper {
	publisher := make(chan publisherData)
	return &LongPollManagerWrapper{
		lpm: initLongPollManager(publisher),
		lpb: initBroker(publisher),
	}
}

func (lpmw *LongPollManagerWrapper) Run() {
	go lpmw.lpb.Run()
	lpmw.lpm.Run()
}

func (lpmw *LongPollManagerWrapper) AddNewToken(token AccessToken) {
	lpmw.lpm.acceptNewTokenCh <- token
}

func (lpmw *LongPollManagerWrapper) DeleteToken(token AccessToken) {
	lpmw.lpm.deleteTokenCh <- token
}

func (lpmw *LongPollManagerWrapper) TokenIsAlive(token AccessToken) <-chan bool {
	ch := make(chan bool)

	lpmw.lpm.isAliveCh <- StatusRequest{token, ch}

	return ch
}

func (lpmw *LongPollManagerWrapper) Subscribe(s *Subscriber) {
	lpmw.lpb.subscribeCh <- s
}

func (lpmw *LongPollManagerWrapper) Unsubscribe(s *Subscriber) {
	lpmw.lpb.unsubscribeCh <- s
}
