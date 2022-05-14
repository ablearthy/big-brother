package longpoll

type Response interface{}

type publisherData struct {
	vkUserId VkUserId
	obj      Response
}

const SUBSCRIBE_ALL = VkUserId(0)

type Subscriber struct {
	VkUserId VkUserId
	Ch       chan Response
}
type longPollBroker struct {
	subscribers   map[*Subscriber]struct{}
	subscribeCh   chan *Subscriber
	unsubscribeCh chan *Subscriber
	publisher     <-chan publisherData
}

func initBroker(publisher <-chan publisherData) *longPollBroker {
	return &longPollBroker{
		subscribers:   make(map[*Subscriber]struct{}),
		subscribeCh:   make(chan *Subscriber),
		unsubscribeCh: make(chan *Subscriber),
		publisher:     publisher,
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
				switch k.VkUserId {
				case SUBSCRIBE_ALL:
					k.Ch <- msg.obj
				case msg.vkUserId:
					k.Ch <- msg.obj
				}
			}
		}
	}
}
