package handler

import (
	"context"

	"github.com/liuzhaomax/go-maxms/internal/core/ext"
	"github.com/liuzhaomax/go-maxms/src/api_user_rpc/code"
	"github.com/liuzhaomax/go-maxms/src/api_user_rpc/model"
	"github.com/liuzhaomax/go-maxms/src/api_user_rpc/pb"
	"github.com/liuzhaomax/go-maxms/src/api_user_rpc/schema"
)

func (h *HandlerUser) GetUserByUserID(ctx context.Context, req *pb.UserIDReq) (*pb.UserRes, error) {
	user := &model.User{}

	err := h.Model.QueryUserByUserID(ctx, req.GetUserID(), user)
	if err != nil {
		h.Logger.Error(ext.FormatError(ext.Unknown, "查询失败", err))

		return nil, code.ErrInternal
	}

	userRes := schema.MapUser2UserRes(user)

	return userRes, nil
}
