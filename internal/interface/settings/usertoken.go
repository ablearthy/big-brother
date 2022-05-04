package settings

type UserTokenSetService interface {
	SetToken(userId int32, accessToken string) error
}

type UserTokenDeleteService interface {
	DeleteToken(userId int32) error
}
