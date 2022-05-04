package settings

import (
	"big-brother/internal/db"
	"big-brother/internal/longpoll"
	"context"
	"database/sql"
	"errors"

	"github.com/SevereCloud/vksdk/v2/api"
)

type UserTokenSetService struct {
	LongPollManager interface {
		AddNewToken(mtr longpoll.ManageTokenRequest)
	}
}

func (utss *UserTokenSetService) SetToken(userId int32, accessToken string) error {
	client := api.NewVK(accessToken)

	me, err := client.UsersGet(api.Params{})
	if err != nil || len(me) != 1 {
		return errors.New("the token is invalid")
	}

	vkUserId := int32(me[0].ID)

	tx, err := db.GetConn().Begin(context.Background())

	if err != nil {
		return errors.New("internal error")
	}

	defer tx.Rollback(context.Background())

	queries := db.New(tx)

	_, err = queries.CreateUserToken(context.Background(), db.CreateUserTokenParams{
		UserID: userId,
		AccessToken: sql.NullString{
			String: accessToken,
			Valid:  true,
		},
	})

	if err != nil {
		return errors.New("internal error")
	}

	_, err = queries.CreateVkToken(context.Background(), db.CreateVkTokenParams{
		AccessToken: accessToken,
		VkUserID:    vkUserId,
	})

	if err != nil {
		return errors.New("internal error")
	}

	tx.Commit(context.Background())
	utss.LongPollManager.AddNewToken(longpoll.ManageTokenRequest{
		AccessToken: longpoll.AccessToken(accessToken),
		VkUserId:    longpoll.VkUserId(vkUserId),
	})
	return nil
}

type UserTokenDeleteService struct {
	LongPollManager interface {
		DeleteToken(mtr longpoll.ManageTokenRequest)
	}
}

func (utds *UserTokenDeleteService) DeleteToken(userId int32) error {

	tx, err := db.GetConn().Begin(context.Background())

	if err != nil {
		return errors.New("internal error")
	}

	defer tx.Rollback(context.Background())

	queries := db.New(tx)

	ut, err := queries.DeleteTokenById(context.Background(), userId)
	if err != nil || !ut.AccessToken.Valid {
		return errors.New("internal error")
	}

	vt, err := queries.GetVkToken(context.Background(), ut.AccessToken.String)

	if err != nil {
		return errors.New("internal error")
	}
	tx.Commit(context.Background())

	utds.LongPollManager.DeleteToken(longpoll.ManageTokenRequest{
		AccessToken: longpoll.AccessToken(ut.AccessToken.String),
		VkUserId:    longpoll.VkUserId(vt.VkUserID),
	})

	return nil
}
