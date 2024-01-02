package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/middleware/auth"
	"github.com/liuzhaomax/go-maxms/internal/middleware/reverse_proxy"
)

var MiddlewareSet = wire.NewSet(wire.Struct(new(Middleware), "*"))

type Middleware struct {
	Auth         *auth.Auth
	ReverseProxy *reverse_proxy.ReverseProxy
}

var MwsSet = wire.NewSet(
	auth.AuthSet,
	reverse_proxy.ReverseProxySet,
)

type IMiddleware interface {
	GenOkMsg(*gin.Context, string) string
	GenErrMsg(*gin.Context, string, error) error
}

var _ IMiddleware = (*auth.Auth)(nil)
var _ IMiddleware = (*reverse_proxy.ReverseProxy)(nil)
