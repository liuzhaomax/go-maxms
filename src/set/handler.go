package set

import (
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/src/data_api/handler"
)

var HandlerSet = wire.NewSet(
	handler.HandlerDataSet,
)
