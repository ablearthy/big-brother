package user

import "big-brother/internal/request/user"

type CreateUserService interface {
	RegisterUser(request *user.CreateUserRequest) error
}
