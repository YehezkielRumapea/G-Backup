package middleware

import (
	"net/http"
	"os"

	"github.com/labstack/echo-jwt/v4" // alias
	"github.com/labstack/echo/v4"
	Echomid "github.com/labstack/echo/v4/middleware"
)

var JwtSecretKey = os.Getenv("JWT_SECRET")

func JwtGuard() echo.MiddlewareFunc {
	if JwtSecretKey == "" {
		JwtSecretKey = "nilainya masih default"
	}

	return Echomid.JWTWithConfig(Echomid.JWTConfig{
		SigningKey:  []byte(JwtSecretKey),
		TokenLookup: "header:Authorization",
		AuthScheme:  "Bearer",
		ErrorHandler: func(c echo.Context, err error) error {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token sudah tidak valid")
		},
	})
}
