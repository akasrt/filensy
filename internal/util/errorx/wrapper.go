package errorx

func WrapUnprocessableEntityError(err error) *Error {
	return &Error{
		Code:    422,
		Message: "unprocessable entity",
		Err:     err,
	}
}

func WrapConflictError(err error) *Error {
	return &Error{
		Code:    409,
		Message: "conflict",
		Err:     err,
	}
}

func WrapInternalServerError(err error) *Error {
	return &Error{
		Code:    500,
		Message: "internal server error",
		Err:     err,
	}
}

func WrapUnauthorizedError(err error, code *string) *Error {
	return &Error{
		Code:    401,
		Message: "unauthorized",
		Err:     err,
		ErrCode: code,
	}
}

func WrapMysqlError(err error) *Error {
	return &Error{
		Err: err,
	}
}

func NewNotFoundError() *Error {
	return &Error{
		Code:    404,
		Message: "not found",
	}
}
