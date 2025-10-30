package middleware

import (
	"log"
	"net/http"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

var jwtSecretKey string

func init() {
	jwtSecretKey = os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		log.Fatal("FATAL: JWT_SECRET_KEY environment variable tidak diatur.")
	}
}

// JWTGuard: Middleware yang melindungi endpoint API.
func JWTGuard() echo.MiddlewareFunc {

	// Gunakan middleware JWT bawaan Echo
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(jwtSecretKey),
		TokenLookup: "header:Authorization", // Mencari token di Authorization header
		ErrorHandler: func(c echo.Context, err error) error {
			// Mengembalikan respons 401 jika token hilang atau tidak valid
			return echo.NewHTTPError(http.StatusUnauthorized, "Token hilang atau tidak valid.")
		},
	})
}
