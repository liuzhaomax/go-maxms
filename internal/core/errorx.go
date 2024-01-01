package core

import (
	"fmt"
	"strconv"
)

// logger.WithField(FAILURE, utils.GetFuncName()).Info(core.FormatError(core.Unknown, "错误描述", err))
// logger.Info(core.FormatInfo("服务启动成功"))

type Code uint32

const (
	OK                Code = 0
	Unknown           Code = 1
	ConfigError       Code = 2
	ConnectionFailed  Code = 3
	ParseIssue        Code = 4
	MissingParameters Code = 400
	Unauthorized      Code = 401
	NotFound          Code = 404
	DownstreamDown    Code = 5
	IOFailure         Code = 6
	PermissionDenied  Code = 7
	CacheDenied       Code = 8
)

func (c Code) String() string {
	switch c {
	case OK:
		return "OK"
	case Unknown:
		return "Unknown"
	case ConfigError:
		return "ConfigError"
	case ConnectionFailed:
		return "ConnectionFailed"
	case ParseIssue:
		return "ParseIssue"
	case MissingParameters:
		return "MissingParameters"
	case Unauthorized:
		return "Unauthorized"
	case NotFound:
		return "NotFound"
	case DownstreamDown:
		return "DownstreamDown"
	case IOFailure:
		return "IOFailure"
	case PermissionDenied:
		return "PermissionDenied"
	case CacheDenied:
		return "CacheDenied"
	default:
		return "Code(" + strconv.FormatInt(int64(c), 10) + ")"
	}
}

type Error struct {
	Code Code // 对错误类型的分类
	Desc string
	Err  error
}

func (err *Error) Error() string {
	if err.Err != nil {
		return fmt.Sprintf("%v: %s: %s", err.Code, err.Desc, err.Err.Error())
	}
	return fmt.Sprintf("%v: %s", err.Code, err.Desc)
}

func FormatInfo(desc string) string {
	return fmt.Sprintf("%v: %s", OK, desc)
}

func FormatError(code Code, desc string, err error) error {
	errObj := new(Error)
	errObj.Code = code
	errObj.Desc = desc
	errObj.Err = err
	return errObj
}
