package user

import (
	"big-brother/internal/db"
	"big-brother/internal/request/user"
	"context"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type LoginUserService struct{}

func comparePasswords(hashedPassword string, notHashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(notHashedPassword))
}

func (_ LoginUserService) Login(request *user.UserLoginRequest) (int, error) {
	ctx := context.Background()
	queries := db.New(db.GetConn())

	hashedPassword, err := hashPassword(request.NotHashedPassword)
	if err != nil {
		return 0, err
	}

	log.Println(request.Username, string(hashedPassword))

	result, err := queries.GetUserByUsername(ctx, request.Username)
	if err != nil {
		return 0, err
	}
	err = comparePasswords(result.Password, request.NotHashedPassword)
	if err != nil {
		return 0, err
	}

	return int(result.ID), nil
}
