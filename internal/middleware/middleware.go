package middleware

import (
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms-template-me/internal/middleware/auth"
)

var MiddlewareSet = wire.NewSet(wire.Struct(new(Middleware), "*"))

type Middleware struct {
	Auth *auth.Auth
}

var MwsSet = wire.NewSet(
	auth.AuthSet,
)
