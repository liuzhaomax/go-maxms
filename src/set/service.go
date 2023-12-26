package set

import (
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms-template/src/data_api/business"
)

var BusinessSet = wire.NewSet(
	business.BusinessDataSet,
)
