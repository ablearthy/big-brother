package longpoll

type Response interface{}

type publisherData struct {
	at  AccessToken
	obj Response
}

type Subscriber struct {
	Token AccessToken
	Ch    chan Response
}
type longPollBroker struct {
	subscribers   map[*Subscriber]struct{}
	subscribeCh   chan *Subscriber
	unsubscribeCh chan *Subscriber
	publisher     <-chan publisherData
}

func initBroker(publisher <-chan publisherData) *longPollBroker {
	return &longPollBroker{
		subscribers: make(map[*Subscriber]struct{}),
		publisher:   publisher,
	}
}

func (lpb *longPollBroker) Subscribe(s *Subscriber) {
	lpb.subscribeCh <- s
}

func (lpb *longPollBroker) Unsubscribe(s *Subscriber) {
	lpb.unsubscribeCh <- s
}

func (lpb *longPollBroker) subscribe(s *Subscriber) {
	lpb.subscribers[s] = struct{}{}
}

func (lpb *longPollBroker) unsubscribe(s *Subscriber) {
	delete(lpb.subscribers, s)
}

func (lpb *longPollBroker) Run() {
	for {
		select {
		case s := <-lpb.subscribeCh:
			lpb.subscribe(s)
		case s := <-lpb.unsubscribeCh:
			lpb.unsubscribe(s)
		case msg := <-lpb.publisher:
			for k := range lpb.subscribers {
				if msg.at == k.Token {
					k.Ch <- msg.obj
				}
			}
		}
	}
}
