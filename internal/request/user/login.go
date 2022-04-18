package user

import (
	"big-brother/internal/utils/validator"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserLoginRequest struct {
	Username          string `json:"username"`
	NotHashedPassword string `json:"password"`
}

type UserLoginRequestValidator struct{}

func (_ *UserLoginRequestValidator) validateUsername(username string) error {
	runeUsername := []rune(username)
	usernameLength := len(runeUsername)
	containsLegalCharacters := true

	for _, v := range runeUsername {
		if !(validator.IsDigit(v) || validator.IsLowerLatin(v)) {
			containsLegalCharacters = false
		}
	}

	if !containsLegalCharacters {
		return errors.New("username contains illegal characters")
	}

	if usernameLength < 4 {
		return errors.New("username is too small")
	}

	if usernameLength > 16 {
		return errors.New("username is too big")
	}

	return nil
}

func (_ *UserLoginRequestValidator) validatePassword(password string) error {
	runePassword := []rune(password)
	passwordLength := len(runePassword)
	containsOnlyLowerLatin := true
	containsIllegalCharacters := false

	for _, v := range runePassword {
		switch {
		case validator.IsPunct(v) || validator.IsUpperLatin(v) || validator.IsDigit(v):
			containsOnlyLowerLatin = false
		case !validator.IsLowerLatin(v):
			containsIllegalCharacters = true
		}
	}
	if containsIllegalCharacters {
		return errors.New("the password contains illegal character")
	}

	if passwordLength > 32 {
		return errors.New("the password is too big")
	}

	if passwordLength < 8 || (containsOnlyLowerLatin && passwordLength < 16) {
		return errors.New("the password is too small")
	}
	return nil
}

func (c *UserLoginRequestValidator) Validate(ulreq *UserLoginRequest) error {
	if err := c.validateUsername(ulreq.Username); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.validatePassword(ulreq.NotHashedPassword); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
