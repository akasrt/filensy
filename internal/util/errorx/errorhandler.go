package errorx

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/akasrt/filensy/internal/util/loggerx"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v5"
)

func ErrorHandler() echo.HTTPErrorHandler {
	return func(c *echo.Context, err error) {
		log := loggerx.Get()

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
			handleMysqlError(appErr)

			if appErr.Code == 500 {
				log.Error("internal error", "requestId", c.Request, "error", appErr.Err.Error())
			}
			if appErr.Code != 0 {
				resp := appErr.Response("")
				c.JSON(appErr.Code, resp)
				return
			}
		}

		log.Error("unknown error", "error", err.Error())
		resp := WrapInternalServerError(err).Response("something unexpected happened")
		c.JSON(500, resp)

	}
}

func handleMysqlError(err *Error) {
	if err.Err == nil {
		return
	}

	if errors.Is(err.Err, sql.ErrNoRows) {
		err.Code = http.StatusNotFound
		err.Message = "request resource was not found"
		errCode := ErrFileNotFound
		err.ErrCode = &errCode
	}

	var mysqlErr *mysql.MySQLError

	if errors.As(err.Err, &mysqlErr) {
		switch mysqlErr.Number {
		case 1062: // duplicate
			err.Code = http.StatusConflict
			err.Message = "conflict"
		case 1452: // fk constraint
			err.Code = http.StatusUnprocessableEntity
			err.Message = "unprocessable entity"
		default:
			err.Code = http.StatusInternalServerError
			err.Message = "internal server error"
		}
	}

}
