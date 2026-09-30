package handler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/zero-hunger/user-service/internal/domain"
)

type Response struct {
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	ResponseData    any    `json:"responseData"`
}

func Success(c echo.Context, status int, message string, data any) error {
	return c.JSON(status, Response{ResponseCode: "00", ResponseMessage: message, ResponseData: data})
}
func Error(c echo.Context, err error) error {
	status, rc, msg := 500, "50", "internal error"
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		status, rc, msg = 400, "40", err.Error()
	case errors.Is(err, domain.ErrUnauthorized):
		status, rc, msg = 401, "41", err.Error()
	case errors.Is(err, domain.ErrNotFound):
		status, rc, msg = 404, "42", err.Error()
	case errors.Is(err, domain.ErrAlreadyExists):
		status, rc, msg = 409, "43", err.Error()
	case errors.Is(err, domain.ErrForbidden):
		status, rc, msg = 403, "44", err.Error()
	}
	return c.JSON(status, Response{ResponseCode: rc, ResponseMessage: msg, ResponseData: nil})
}
