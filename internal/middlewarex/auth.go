package middlewarex

import (
	"net/http"
	"strings"

	"github.com/akasrt/filensy/internal/config/env"
	"github.com/labstack/echo/v5"
)

func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authKey := env.GetEnv(env.AuthKey)
			if strings.TrimSpace(authKey) == "" {
				return next(c)
			}

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			const prefix = "Bearer "

			if !strings.HasPrefix(authHeader, prefix) {
				return echo.NewHTTPError(http.StatusUnauthorized, "authentication failed")
			}

			token := strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))

			if token != authKey {
				return echo.NewHTTPError(http.StatusUnauthorized, "authentication failed")
			}

			return next(c)
		}
	}
}
