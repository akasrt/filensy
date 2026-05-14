package httputil

type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   *Error `json:"error,omitempty"`
}

type Error struct {
	Code     int    `json:"code,omitempty"`
	ErrorMsg string `json:"error_message,omitempty"`
}

func NewResponse(message string, data any) Response {
	return Response{
		Message: message,
		Data:    data,
	}
}

func NewErrorResponse(message string, err *Error) Response {
	return Response{
		Message: message,
		Error:   err,
	}
}
