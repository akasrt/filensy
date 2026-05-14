package errorx

func NewUnprocessableEntityError(err error) *Error {
	return &Error{
		Code:    422,
		Message: "unprocessable entity",
		Err:     err,
	}
}

func NewConflictError(err error) *Error {
	return &Error{
		Code:    409,
		Message: "conflict",
		Err:     err,
	}
}

func NewInternalServerError(err error) *Error {
	return &Error{
		Code:    500,
		Message: "internal server error",
		Err:     err,
	}
}

func NewUnauthorizedError(err error) *Error {
	return &Error{
		Code:    401,
		Message: "unauthorized",
		Err:     err,
	}
}
