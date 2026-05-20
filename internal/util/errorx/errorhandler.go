package errorx

import (
	"errors"

	"github.com/akasrt/filensy/internal/util/loggerx"
	"github.com/labstack/echo/v5"
)

func ErrorHandler() echo.HTTPErrorHandler {
	return func(c *echo.Context, err error) {
		var sc echo.HTTPStatusCoder
		if errors.As(err, &sc) {
			statusNotFound := 404
			if sc.StatusCode() == statusNotFound {
				c.JSON(statusNotFound, NewNotFoundError().Response("resource not found"))
				return
			}
		}
		var appErr *Error
		if errors.As(err, &appErr) {
			resp := appErr.Response("")
			c.JSON(appErr.Code, resp)
		} else {
			log := loggerx.Get()
			log.Error("unknown error", "error", err.Error())
			resp := NewInternalServerError(err).Response("something unexpected happened")
			c.JSON(500, resp)
		}
	}
}
