package code

type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Error(code int, msg string) ApiError {
	return ApiError{
		Code:    code,
		Message: msg,
	}
}

func (code ApiError) Error() string {
	return code.Message
}

var (
	InternalErr = Error(5001, "internal error")
)
