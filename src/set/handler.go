package set

import (
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/src/api_user/handler"
	handlerRpc "github.com/liuzhaomax/go-maxms/src/api_user_rpc/handler"
)

var HandlerSet = wire.NewSet(
	handler.HandlerUserSet,
	handlerRpc.HandlerUserSet,
)
