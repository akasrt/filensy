package middlewarex

import (
	"github.com/akasrt/filensy/internal/util/loggerx"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func LoggerMiddleware() echo.MiddlewareFunc {
	log := loggerx.Get()

	return middleware.RequestLoggerWithConfig(
		middleware.RequestLoggerConfig{
			LogMethod:   true,
			LogURI:      true,
			LogStatus:   true,
			LogLatency:  true,
			LogRemoteIP: true,
			LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
				reqID := c.Response().Header().Get(echo.HeaderXRequestID)
				args := []any{
					"request_id", reqID,
					"method", v.Method,
					"path", v.URI,
					"status", v.Status,
					"latency", v.Latency.String(),
					"ip", v.RemoteIP,
				}

				switch {
				case v.Status >= 500:
					log.Error("http request", args...)
				case v.Status >= 400:
					log.Warn("http request", args...)
				default:
					log.Info("http request", args...)
				}

				return nil
			},
		},
	)
}
