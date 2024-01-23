package business

import (
	"context"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/liuzhaomax/go-maxms/src/api_user_rpc/model"
	"github.com/liuzhaomax/go-maxms/src/api_user_rpc/pb"
	"github.com/liuzhaomax/go-maxms/src/api_user_rpc/schema"
	"github.com/redis/go-redis/v9"
)

var BusinessUserSet = wire.NewSet(wire.Struct(new(BusinessUser), "*"))

type BusinessUser struct {
	Model *model.ModelUser
	Tx    *core.Trans
	Redis *redis.Client
}

func (b *BusinessUser) GetUserByUserID(ctx context.Context, req *pb.UserIDReq) (*pb.UserRes, error) {
	user := &model.User{}
	err := b.Model.QueryUserByUserID(ctx, req.UserID, user)
	if err != nil {
		return nil, err
	}
	userRes := schema.MapUser2UserRes(user)
	return userRes, nil
}
