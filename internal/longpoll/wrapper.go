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

func (lpmw *LongPollManagerWrapper) AddNewToken(mtr ManageTokenRequest) {
	lpmw.lpm.acceptNewTokenCh <- mtr
}

func (lpmw *LongPollManagerWrapper) DeleteToken(mtr ManageTokenRequest) {
	lpmw.lpm.deleteTokenCh <- mtr
}

func (lpmw *LongPollManagerWrapper) TokenIsAlive(vkUserId VkUserId) <-chan bool {
	ch := make(chan bool)

	lpmw.lpm.isAliveCh <- StatusRequest{vkUserId, ch}

	return ch
}

func (lpmw *LongPollManagerWrapper) Subscribe(s *Subscriber) {
	lpmw.lpb.subscribeCh <- s
}

func (lpmw *LongPollManagerWrapper) Unsubscribe(s *Subscriber) {
	lpmw.lpb.unsubscribeCh <- s
}
