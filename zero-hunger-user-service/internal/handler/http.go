package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"github.com/zero-hunger/user-service/internal/domain"
	"github.com/zero-hunger/user-service/internal/service"
)

type HTTPHandler struct{ users *service.UserService }

func NewHTTPHandler(users *service.UserService) *HTTPHandler { return &HTTPHandler{users: users} }

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
type profileRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func (h *HTTPHandler) Register(c echo.Context) error {
	var in registerRequest
	if c.Bind(&in) != nil {
		return Error(c, domain.ErrInvalidInput)
	}
	u, err := h.users.Register(c.Request().Context(), service.RegisterInput{Name: in.Name, Email: in.Email, Phone: in.Phone, Password: in.Password, Role: domain.Role(in.Role)})
	if err != nil {
		return Error(c, err)
	}
	return Success(c, http.StatusCreated, "user registered successfully", u)
}
func (h *HTTPHandler) Login(c echo.Context) error {
	var in loginRequest
	if c.Bind(&in) != nil {
		return Error(c, domain.ErrInvalidInput)
	}
	out, err := h.users.Login(c.Request().Context(), in.Email, in.Password)
	if err != nil {
		return Error(c, err)
	}
	return Success(c, http.StatusOK, "login successfully", out)
}
func (h *HTTPHandler) Refresh(c echo.Context) error {
	var in refreshRequest
	if c.Bind(&in) != nil {
		return Error(c, domain.ErrInvalidInput)
	}
	out, err := h.users.Refresh(c.Request().Context(), in.RefreshToken)
	if err != nil {
		return Error(c, err)
	}
	return Success(c, http.StatusOK, "token refreshed successfully", out)
}
func (h *HTTPHandler) Logout(c echo.Context) error {
	var in refreshRequest
	if c.Bind(&in) != nil {
		return Error(c, domain.ErrInvalidInput)
	}
	if err := h.users.Logout(c.Request().Context(), in.RefreshToken); err != nil {
		return Error(c, err)
	}
	return Success(c, http.StatusOK, "logout successfully", nil)
}
func (h *HTTPHandler) Me(c echo.Context) error {
	id, _, err := h.authID(c)
	if err != nil {
		return Error(c, err)
	}
	u, err := h.users.GetByID(c.Request().Context(), id)
	if err != nil {
		return Error(c, err)
	}
	return Success(c, 200, "profile retrieved successfully", u)
}
func (h *HTTPHandler) UpdateMe(c echo.Context) error {
	id, _, err := h.authID(c)
	if err != nil {
		return Error(c, err)
	}
	var in profileRequest
	if c.Bind(&in) != nil {
		return Error(c, domain.ErrInvalidInput)
	}
	u, err := h.users.UpdateProfile(c.Request().Context(), id, in.Name, in.Phone)
	if err != nil {
		return Error(c, err)
	}
	return Success(c, 200, "profile updated successfully", u)
}
func (h *HTTPHandler) GetUser(c echo.Context) error {
	id, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		return Error(c, domain.ErrInvalidInput)
	}
	u, err := h.users.GetByID(c.Request().Context(), id)
	if err != nil {
		return Error(c, err)
	}
	return Success(c, 200, "user retrieved successfully", u)
}
func (h *HTTPHandler) authID(c echo.Context) (uuid.UUID, domain.Role, error) {
	header := c.Request().Header.Get("Authorization")
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return uuid.Nil, "", domain.ErrUnauthorized
	}
	return h.users.ParseAccessToken(parts[1])
}
