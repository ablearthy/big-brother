package postinit

import (
	"big-brother/internal/background"
	"big-brother/internal/db"
	"big-brother/internal/longpoll"
	"context"
)

func StartLongPollForAllUsers() error {
	queries := db.New(db.GetConn())

	users, err := queries.GetAllUserTokens(context.Background())
	if err != nil {
		return err
	}

	lpmw := background.GetLongPollManagerWrapper()

	for _, u := range users {
		if !u.VkUserID.Valid || !u.AccessToken.Valid {
			continue
		}
		lpmw.AddNewToken(longpoll.ManageTokenRequest{
			AccessToken: longpoll.AccessToken(u.AccessToken.String),
			VkUserId:    longpoll.VkUserId(u.VkUserID.Int32),
		})
	}

	return nil
}
