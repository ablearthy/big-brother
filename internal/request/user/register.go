package user

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CreateUserRequest struct {
	Username          string `json:"username"`
	NotHashedPassword string `json:"password"`
	InviteCode        string `json:"invite_code"`
}

type CreateUserRequestValidator struct{}

func isPunct(ch rune) bool {
	return ('!' <= ch && ch <= '/') || (':' <= ch && ch <= '@') || ('[' <= ch && ch <= '`') || ('{' <= ch && ch <= '~')
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func isUpperLatin(ch rune) bool {
	return 'A' <= ch && ch <= 'Z'
}

func isLowerLatin(ch rune) bool {
	return 'a' <= ch && ch <= 'z'
}

func (_ *CreateUserRequestValidator) validateUsername(username string) error {
	runeUsername := []rune(username)
	usernameLength := len(runeUsername)
	containsLegalCharacters := true

	for _, v := range runeUsername {
		if !(isDigit(v) || isLowerLatin(v)) {
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

func (_ *CreateUserRequestValidator) validatePassword(password string) error {
	runePassword := []rune(password)
	passwordLength := len(runePassword)
	containsOnlyLowerLatin := true
	containsIllegalCharacters := false

	for _, v := range runePassword {
		switch {
		case isPunct(v) || isUpperLatin(v) || isDigit(v):
			containsOnlyLowerLatin = false
		case !isLowerLatin(v):
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

func (_ *CreateUserRequestValidator) validateInviteCode(inviteCode string) error {
	runeInviteCode := []rune(inviteCode)
	validInviteCode := true

	for _, v := range runeInviteCode {
		if !(isLowerLatin(v) || isUpperLatin(v)) {
			validInviteCode = false
		}
	}
	if !validInviteCode || len(runeInviteCode) != 10 {
		return errors.New("the invite code is invalid")
	}
	return nil
}

func (c *CreateUserRequestValidator) Validate(cureq *CreateUserRequest) error {
	if err := c.validateUsername(cureq.Username); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.validatePassword(cureq.NotHashedPassword); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.validateInviteCode(cureq.InviteCode); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
