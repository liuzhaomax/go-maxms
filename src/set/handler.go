package set

import (
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms-template/src/data_api/handler"
)

var HandlerSet = wire.NewSet(
	handler.HandlerDataSet,
)
