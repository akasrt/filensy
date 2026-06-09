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
	ErrInvalidConfigKey    error = errors.New("invalid config key")
)

// CLIErrorHandler centralized error handler for cli
func CLIErrorHandler(err error) (exit int, message string) {
	if err == nil {
		return 0, "success"
	}

	if errors.Is(err, ErrLocalCreationFailed) || errors.Is(err, ErrLocalDeletionFailed) {
		return 0, fmt.Sprintf("Warning: %v", err)
	}

	if errors.Is(err, ErrFileNotExists) {
		return 1, "Error: The specified file does not exist."
	}

	if errors.Is(err, ErrPasswordMissing) {
		return 1, "Error: A password is required for this encrypted file."
	}

	if errors.Is(err, ErrInvalidConfigKey) {
		return 1, "Invalid config key! Supported keys are: 'auth', 'dir'"
	}

	var serverErr *ServerError
	if errors.As(err, &serverErr) {
		return 1, fmt.Sprintf("Server Error [%d]: %s", serverErr.Status, serverErr.Resp.Message)
	}

	return 1, fmt.Sprintf("An unexpected error occurred: %v", err)
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
