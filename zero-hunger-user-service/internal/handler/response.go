package handler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/zero-hunger/user-service/internal/domain"
)

type Response struct {
	RC      string `json:"rc"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Success(c echo.Context, status int, message string, data any) error {
	return c.JSON(status, Response{RC: "00", Message: message, Data: data})
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
	return c.JSON(status, Response{RC: rc, Message: msg, Data: nil})
}
