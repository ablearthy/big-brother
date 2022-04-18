package user

import (
	"big-brother/internal/request/user"
)

type LoginUserService interface {
	Login(request *user.UserLoginRequest) (int, error)
}
