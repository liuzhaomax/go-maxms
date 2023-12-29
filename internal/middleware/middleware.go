package middleware

import (
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
