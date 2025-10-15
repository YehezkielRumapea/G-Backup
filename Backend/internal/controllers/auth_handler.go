package controllers

import (
	"gbackup-system/backend/internal/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	AuthSvc services.AuthService
}

func NewAuthHandler(e *echo.Echo, authSvc services.AuthService) *AuthHandler {
	return &AuthHandler{AuthSvc: authSvc}
}

// Endpoint Login
func (h *AuthHandler) Login(c echo.Context) error {
	req := new(services.LoginRequest)

	// Json to struct
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"Error": "Format Input Salah"})
	}

	// AuthService untuk Generate Token
	token, err := h.AuthSvc.Authenticate(req)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"Error": "Kredensial Salah"})
	}

	// Token JWt ke frontend
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Login Berhasil",
		"token":   token,
	})
}

// func (h *AuthHandler) RegisterAdmin(c.echo.echo.Context) error {

// }
