package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/middleware/auth"
	"github.com/liuzhaomax/go-maxms/internal/middleware/reverse_proxy"
	"github.com/liuzhaomax/go-maxms/internal/middleware/tracing"
	"github.com/liuzhaomax/go-maxms/internal/middleware/validator"
	"github.com/liuzhaomax/go-maxms/internal/middleware/ws_upgrader"
)

var MiddlewareSet = wire.NewSet(wire.Struct(new(Middleware), "*"))

type Middleware struct {
	Auth         *auth.Auth
	Validator    *validator.Validator
	Tracing      *tracing.Tracing
	ReverseProxy *reverse_proxy.ReverseProxy
	wsUpgrader   *ws_upgrader.WsUpgrader
}

var MwsSet = wire.NewSet(
	auth.AuthSet,
	validator.ValidatorSet,
	tracing.TracingSet,
	reverse_proxy.ReverseProxySet,
	ws_upgrader.WsUpgraderSet,
)

type IMiddleware interface {
	AbortWithError(*gin.Context, ...any)
}

var _ IMiddleware = (*auth.Auth)(nil)
var _ IMiddleware = (*validator.Validator)(nil)
var _ IMiddleware = (*tracing.Tracing)(nil)
var _ IMiddleware = (*reverse_proxy.ReverseProxy)(nil)
var _ IMiddleware = (*ws_upgrader.WsUpgrader)(nil)
