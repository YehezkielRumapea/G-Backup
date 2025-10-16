package middleware

import (
	"net/http"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4" // alias
	"github.com/labstack/echo/v4"
)

var JwtSecretKey = os.Getenv("JWT_SECRET")

func JwtGuard() echo.MiddlewareFunc {
	if JwtSecretKey == "" {
		JwtSecretKey = "nilainya masih default"
	}

	return echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(JwtSecretKey),
		TokenLookup: "header:Authorization",
		ErrorHandler: func(c echo.Context, err error) error {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token sudah tidak valid")
		},
	})
}
