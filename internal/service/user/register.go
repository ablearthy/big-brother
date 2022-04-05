package user

import (
	"big-brother/internal/db"
	"big-brother/internal/request/user"
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
)

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 14)
}

func createInviteCode(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	s := make([]byte, length)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

type CreateUserService struct{}

var createUserError = fmt.Errorf("an error occurred while creating user")

func (_ CreateUserService) RegisterUser(request *user.CreateUserRequest) error {
	ctx := context.Background()
	queries := db.New(db.GetConn())

	inviterId, err := queries.GetUserIdByInviteCode(ctx, request.InviteCode)
	if err != nil {
		return createUserError
	}

	cntUsed, err := queries.GetCountOfUsedInviteCodes(context.Background(), inviterId)
	if err != nil || cntUsed >= 5 {
		return createUserError
	}

	tx, err := db.GetConn().Begin(context.Background())
	if err != nil {
		return createUserError
	}
	defer tx.Rollback(context.Background())

	queries = queries.WithTx(tx)

	hashedPassword, err := hashPassword(request.NotHashedPassword)

	if err != nil {
		return createUserError
	}

	newInviteCode := createInviteCode(10)

	// NOTE: Since the username field in database is unique, there's no need to check
	// if username exists in the table
	u, err := queries.CreateUser(context.Background(), db.CreateUserParams{
		Username:  request.Username,
		Password:  string(hashedPassword),
		InviterID: inviterId,
	})

	if err != nil {
		return createUserError
	}
	_, err = queries.CreateInviteCode(context.Background(), db.CreateInviteCodeParams{
		UserID:     u.ID,
		InviteCode: newInviteCode,
	})

	if err != nil {
		return createUserError
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return createUserError
	}
	return nil
}
