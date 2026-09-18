package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zero-hunger/request-service/internal/domain"
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
	status, rc, message := http.StatusInternalServerError, "50", "internal error"
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		status, rc, message = 400, "40", err.Error()
	case errors.Is(err, domain.ErrUnauthorized):
		status, rc, message = 401, "41", err.Error()
	case errors.Is(err, domain.ErrNotFound):
		status, rc, message = 404, "42", err.Error()
	case errors.Is(err, domain.ErrConflict):
		status, rc, message = 409, "43", err.Error()
	case errors.Is(err, domain.ErrForbidden):
		status, rc, message = 403, "44", err.Error()
	}
	return c.JSON(status, Response{RC: rc, Message: message, Data: nil})
}
