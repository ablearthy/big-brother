package background

import (
	"big-brother/internal/longpoll/dbsave"
)

var dblps *dbsave.DbLongPollSaver

func InitDbLongPollSaver() {
	dblps = dbsave.New()
}

func GetDbLongPollSaver() *dbsave.DbLongPollSaver {
	return dblps
}
