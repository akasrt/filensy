package errorx

import (
	"fmt"

	"github.com/akasrt/filensy/internal/util/httputil"
)

type Error struct {
	Code    int
	Message string
	Err     error
	ErrCode *string
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprint(e.Message + " " + e.Err.Error())
	} else {
		return fmt.Sprint(e.Message)
	}
}

func (e *Error) Unwrap() error {
	return e.Err
}

func New(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func Wrap(err error, code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (e *Error) Response(message string) httputil.Response {
	return httputil.NewErrorResponse(message, &httputil.Error{
		Code:     e.ErrCode,
		ErrorMsg: e.Message,
	})
}
