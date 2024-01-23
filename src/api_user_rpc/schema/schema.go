package schema

import (
	"github.com/liuzhaomax/go-maxms/src/api_user_rpc/model"
	"github.com/liuzhaomax/go-maxms/src/api_user_rpc/pb"
)

func MapUser2UserRes(user *model.User) *pb.UserRes {
	return &pb.UserRes{
		Id:     int32(user.ID),
		Mobile: user.Mobile,
	}
}
