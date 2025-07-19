package code

import "github.com/liuzhaomax/go-maxms/internal/core"

func Error(code int, msg string) *core.ApiError {
	return &core.ApiError{
		Code:    code,
		Message: msg,
	}
}

var (
	// 内部服务器错误 1000-1999
	InternalErr = Error(1000, "内部错误")
	DBFailed    = Error(1001, "数据库错误")
	// 数据错误 2000-2999
	// 下游服务错误 10000+
)
