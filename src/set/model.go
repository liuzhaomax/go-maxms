package set

import (
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/src/api_user/model"
	modelRpc "github.com/liuzhaomax/go-maxms/src/api_user_rpc/model"
)

var ModelSet = wire.NewSet(
	model.ModelUserSet,
	modelRpc.ModelUserSet,
)
