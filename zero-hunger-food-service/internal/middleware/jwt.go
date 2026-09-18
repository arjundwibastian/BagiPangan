package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/zero-hunger/food-service/internal/dto"
)

type JWTMiddleware struct {
	jwtSecret []byte
}

func NewJWTMiddleware(jwtSecret []byte) *JWTMiddleware {
	return &JWTMiddleware{jwtSecret: jwtSecret}
}

func (m *JWTMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: "Error: Token not found",
				ResponseData:    nil,
			})
		}

		tokenString := ""
		splitString := strings.Split(authHeader, " ")
		if len(splitString) > 1 && splitString[0] == "Bearer" {
			tokenString = splitString[1]
		} else {
			return c.JSON(http.StatusUnauthorized, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: "Error: Token not found",
				ResponseData:    nil,
			})
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return m.jwtSecret, nil
		})
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: "Error: Token is invalid",
				ResponseData:    nil,
			})
		}

		claims := token.Claims.(jwt.MapClaims)
		userID, ok := claims["sub"].(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: "Error: Token is invalid",
				ResponseData:    nil,
			})
		}
		if _, err := uuid.Parse(userID); err != nil {
			return c.JSON(http.StatusUnauthorized, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: "Error: user id is invalid",
				ResponseData:    nil,
			})
		}
		role, ok := claims["role"].(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: "Error: role is invalid",
				ResponseData:    nil,
			})
		}
		c.Set("user_id", userID)
		c.Set("role", role)

		return next(c)
	}
}

func (m *JWTMiddleware) DonorOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, ok := c.Get("role").(string)
		if !ok || role != "donor" {
			return c.JSON(http.StatusForbidden, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: "Error: donor user access required",
			})
		}
		return next(c)
	}
}
