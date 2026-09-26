package handler

import (
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/zero-hunger/request-service/internal/domain"
	"github.com/zero-hunger/request-service/internal/usecase"
)

type HTTPHandler struct {
	requests *usecase.RequestService
	secret   string
}

func NewHTTPHandler(requests *usecase.RequestService, secret string) *HTTPHandler {
	return &HTTPHandler{requests: requests, secret: secret}
}

type createRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	RadiusKM  float64 `json:"radius_km"`
}

func (h *HTTPHandler) Create(c echo.Context) error {
	userID, role, err := h.auth(c)
	if err != nil {
		return Error(c, err)
	}
	if role != "recipient" {
		return Error(c, domain.ErrForbidden)
	}
	var in createRequest
	if c.Bind(&in) != nil {
		return Error(c, domain.ErrInvalidInput)
	}
	out, err := h.requests.Create(c.Request().Context(), usecase.CreateInput{UserID: userID, Latitude: in.Latitude, Longitude: in.Longitude, RadiusKM: in.RadiusKM})
	if err != nil {
		return Error(c, err)
	}
	return Success(c, 201, "food request created successfully", out)
}

func (h *HTTPHandler) Get(c echo.Context) error {
	userID, role, err := h.auth(c)
	if err != nil {
		return Error(c, err)
	}
	id, err := uuid.Parse(c.Param("request_id"))
	if err != nil {
		return Error(c, domain.ErrInvalidInput)
	}
	request, err := h.requests.Get(c.Request().Context(), id)
	if err != nil {
		return Error(c, err)
	}
	if role != "admin" && request.UserID != userID {
		return Error(c, domain.ErrForbidden)
	}
	return Success(c, 200, "food request retrieved successfully", request)
}

func (h *HTTPHandler) ListByUser(c echo.Context) error {
	userID, role, err := h.auth(c)
	if err != nil {
		return Error(c, err)
	}
	targetID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		return Error(c, domain.ErrInvalidInput)
	}
	if role != "admin" && targetID != userID {
		return Error(c, domain.ErrForbidden)
	}
	requests, err := h.requests.ListByUser(c.Request().Context(), targetID)
	if err != nil {
		return Error(c, err)
	}
	return Success(c, 200, "food requests retrieved successfully", requests)
}

func (h *HTTPHandler) Cancel(c echo.Context) error {
	userID, _, err := h.auth(c)
	if err != nil {
		return Error(c, err)
	}
	id, err := uuid.Parse(c.Param("request_id"))
	if err != nil {
		return Error(c, domain.ErrInvalidInput)
	}
	request, err := h.requests.Cancel(c.Request().Context(), id, userID)
	if err != nil {
		return Error(c, err)
	}
	return Success(c, 200, "food request cancelled successfully", request)
}

func (h *HTTPHandler) auth(c echo.Context) (uuid.UUID, string, error) {
	parts := strings.SplitN(c.Request().Header.Get("Authorization"), " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return uuid.Nil, "", domain.ErrUnauthorized
	}
	token, err := jwt.Parse(parts[1], func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrUnauthorized
		}
		return []byte(h.secret), nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil, "", domain.ErrUnauthorized
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, "", domain.ErrUnauthorized
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, "", domain.ErrUnauthorized
	}
	id, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, "", domain.ErrUnauthorized
	}
	role, _ := claims["role"].(string)
	return id, role, nil
}
