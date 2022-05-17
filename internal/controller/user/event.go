package user

import (
	"big-brother/internal/auth"
	svc "big-brother/internal/interface/user"
	req "big-brother/internal/request/user"
	"net/http"

	"github.com/labstack/echo/v4"
)

type LastMessageEventHandler struct {
	Validator *req.LastMessageEventRequestValidator
	Service   svc.LastMessageEventService
}

func (lmeh *LastMessageEventHandler) GetMessageEvents(c echo.Context) error {

	userId, err := auth.GetUserId(c)

	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	lmereq := new(req.LastMessageEventRequest)
	if err := c.Bind(lmereq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := lmeh.Validator.Validate(lmereq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	data, err := lmeh.Service.GetMessageEvents(userId, lmereq.LastId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, data)
}
