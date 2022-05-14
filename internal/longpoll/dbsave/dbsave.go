package dbsave

import (
	"big-brother/internal/longpoll"
	"log"
)

type DbLongPollSaver struct {
	ch chan longpoll.Response
}

func New() *DbLongPollSaver {
	return &DbLongPollSaver{
		ch: make(chan longpoll.Response),
	}
}

func (dblps *DbLongPollSaver) GetChannel() chan longpoll.Response {
	return dblps.ch
}

func (dblps *DbLongPollSaver) Run() {
	for {
		data := <-dblps.ch
		log.Println("DBSAVE:", data)
	}
}
