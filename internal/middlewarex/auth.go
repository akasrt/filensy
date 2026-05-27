package middlewarex

import (
	"strings"

	"github.com/akasrt/filensy/internal/config/env"
	"github.com/akasrt/filensy/internal/util/errorx"
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
				errCode := errorx.ErrMissingAuthToken
				return errorx.WrapUnauthorizedError(nil, &errCode)
			}

			const prefix = "Bearer "

			if !strings.HasPrefix(authHeader, prefix) {
				errCode := errorx.ErrInvalidAuthToken
				return errorx.WrapUnauthorizedError(nil, &errCode)
			}

			token := strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))

			if token != authKey {
				errCode := errorx.ErrInvalidAuthToken
				return errorx.WrapUnauthorizedError(nil, &errCode)
			}

			return next(c)
		}
	}
}
