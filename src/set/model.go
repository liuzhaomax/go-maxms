package set

import (
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/src/api_user/model"
)

var ModelSet = wire.NewSet(
	model.ModelUserSet,
)
