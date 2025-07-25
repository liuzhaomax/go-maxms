package ext

import (
	"fmt"
	"strconv"
)

// logger.WithField(FAILURE, utils.GetFuncName()).Info(core.FormatError(core.Unknown, "错误描述", err))
// logger.Info(core.FormatInfo("服务启动成功"))

type Code uint32

const (
	OK                    Code = 0
	Unknown               Code = 1
	ConfigError           Code = 2
	ConnectionFailed      Code = 3
	ParseIssue            Code = 4
	MissingParameters     Code = 400
	Unauthorized          Code = 401
	Forbidden             Code = 403
	NotFound              Code = 404
	InternalServerError   Code = 500
	DownstreamDown        Code = 5
	IOException           Code = 6
	PermissionDenied      Code = 7
	DBDenied              Code = 8
	CacheDenied           Code = 9
	VaultDenied           Code = 10
	ProtocolUpgradeFailed Code = 11
	CloseException        Code = 12
	CommunicationFailed   Code = 13
	TypeError             Code = 14
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
	case ProtocolUpgradeFailed:
		return "协议升级失败"
	case CloseException:
		return "关闭异常"
	case CommunicationFailed:
		return "通信异常"
	case TypeError:
		return "类型错误"
	default:
		return "Code(" + strconv.FormatInt(int64(c), 10) + ")"
	}
}

type Error struct {
	Code Code   `json:"code"`
	Desc string `json:"desc"`
}

func (err *Error) Error() string {
	return fmt.Sprintf("%v: %s", err.Code, err.Desc)
}

func FormatInfo(desc string) string {
	return fmt.Sprintf("%v: %s", OK, desc)
}

func FormatError(code Code, desc string, err error) error {
	errObj := new(Error)
	errObj.Code = code
	errObj.Desc = desc
	if err != nil {
		errObj.Desc = fmt.Sprintf("%s: %s", desc, err.Error())
	}
	return errObj
}

func FormatCaller(ok bool, desc string) string {
	if ok {
		return "成功: Caller: " + desc
	}

	return "失败: Caller: " + desc
}

type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (code *ApiError) Error() string {
	return code.Message
}
