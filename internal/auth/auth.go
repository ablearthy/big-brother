package auth

import (
	"errors"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const ONE_DAY = 86400

func SetCookies(c echo.Context, userId int) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   7 * ONE_DAY,
		HttpOnly: true,
	}
	sess.Values["userId"] = userId
	err = sess.Save(c.Request(), c.Response())
	return err
}

func RefreshCookies(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   7 * ONE_DAY,
		HttpOnly: true,
	}
	if userId, ok := sess.Values["userId"]; ok {
		sess.Values["userId"] = userId
	}
	err = sess.Save(c.Request(), c.Response())
	return err
}

func GetUserId(c echo.Context) (int, error) {
	sess, err := session.Get("session", c)
	if err != nil {
		return 0, err
	}
	if userId, ok := sess.Values["userId"].(int); ok {
		return userId, nil
	}
	return 0, errors.New("userId does not exist")
}
