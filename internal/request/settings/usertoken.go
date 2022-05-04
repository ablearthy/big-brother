package settings

import "errors"

type UserTokenSetRequest struct {
	AccessToken string `json:"access_token"`
}

type UserTokenSetRequestValidator struct{}

func (*UserTokenSetRequestValidator) Validate(utsr *UserTokenSetRequest) error {
	token := []rune(utsr.AccessToken)
	if len(token) > 100 {
		return errors.New("access_token is too long")
	}
	return nil
}
