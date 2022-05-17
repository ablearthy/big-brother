package user

import "errors"

type LastMessageEventRequest struct {
	LastId int `json:"last_id"`
}

type LastMessageEventRequestValidator struct{}

func (*LastMessageEventRequestValidator) Validate(req *LastMessageEventRequest) error {
	if req.LastId < 0 {
		return errors.New("last_id must not be non-negative")
	}
	return nil
}
