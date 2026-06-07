package errorx

import (
	"errors"
	"fmt"

	"github.com/akasrt/filensy/internal/util/httputil"
)

var (
	ErrFileNotExists       error = errors.New("file doesn't exists")
	ErrPasswordMissing     error = errors.New("password is needed for encrypted files")
	ErrLocalDeletionFailed error = errors.New("unable to delete local file data")
	ErrLocalCreationFailed error = errors.New("unable to create local file data")
)

// CLIErrorHandler centralized error handler for cli
func CLIErrorHandler(err error) (exit int, message string) {
	return 0, "success"
}

func WrapServerError(status int, resp httputil.Response) error {
	return &ServerError{Status: status, Resp: resp}
}

type ServerError struct {
	Status int
	Resp   httputil.Response
}

func (se *ServerError) Error() string {
	return fmt.Sprintf("server error occurred! Status: %d, Message: %s", se.Status, se.Resp.Message)
}
