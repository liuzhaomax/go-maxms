package core

import (
	"fmt"
	"strconv"
)

// logger.WithField(FAILURE, utils.GetFuncName()).Info(core.FormatError(core.Unknown, "错误描述", err))
// logger.Info(core.FormatInfo("服务启动成功"))

type Code uint32

const (
	OK                  Code = 0
	Unknown             Code = 1
	ConfigError         Code = 2
	ConnectionFailed    Code = 3
	ParseIssue          Code = 4
	MissingParameters   Code = 400
	Unauthorized        Code = 401
	Forbidden           Code = 403
	NotFound            Code = 404
	InternalServerError Code = 500
	DownstreamDown      Code = 5
	IOException         Code = 6
	PermissionDenied    Code = 7
	DBDenied            Code = 8
	CacheDenied         Code = 9
	VaultDenied         Code = 10
)

func (c Code) String() string {
	switch c {
	case OK:
		return "OK"
	case Unknown:
		return "Unknown"
	case ConfigError:
		return "配置错误"
	case ConnectionFailed:
		return "连接失败"
	case ParseIssue:
		return "解析问题"
	case MissingParameters:
		return "缺少参数"
	case Unauthorized:
		return "未授权"
	case Forbidden:
		return "请求被拒绝"
	case NotFound:
		return "没找到"
	case InternalServerError:
		return "内部服务器错误"
	case DownstreamDown:
		return "下游宕机"
	case IOException:
		return "IO异常"
	case PermissionDenied:
		return "无权限"
	case DBDenied:
		return "数据库拒绝"
	case CacheDenied:
		return "缓存拒绝"
	case VaultDenied:
		return "Vault拒绝"
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

func FormatCaller(ok bool, desc string) string {
	if ok {
		return fmt.Sprintf("%s: Caller: %s", SUCCESS, desc)
	}
	return fmt.Sprintf("%s: Caller: %s", FAILURE, desc)
}
