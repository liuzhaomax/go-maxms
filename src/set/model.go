package set

import (
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms-template/src/data_api/model"
)

var ModelSet = wire.NewSet(
	model.ModelDataSet,
)
