package set

import (
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/src/api_user/business"
)

var BusinessSet = wire.NewSet(
	business.BusinessUserSet,
)
