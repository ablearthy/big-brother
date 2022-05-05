package background

import "big-brother/internal/longpoll"

var lpmw *longpoll.LongPollManagerWrapper

func InitLongPollManagerWrapper() {
	lpmw = longpoll.Init()
}

func GetLongPollManagerWrapper() *longpoll.LongPollManagerWrapper {
	return lpmw
}
