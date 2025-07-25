package validator

import (
	"context"

	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core/config"
	"github.com/liuzhaomax/go-maxms/internal/core/ext"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var ValidatorRPCSet = wire.NewSet(wire.Struct(new(ValidatorRPC), "*"))

type ValidatorRPC struct {
	Logger *logrus.Logger
	Redis  *redis.Client
}

func (v *ValidatorRPC) ValidateMetadata(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		v.AbortWithError(ctx, ext.ParseIssue, "元数据解析错误", err)

		return
	}

	err = config.ValidateMetadata(md)
	if err != nil {
		v.AbortWithError(ctx, ext.ParseIssue, "元数据错误", err)

		return
	}

	resp, err = handler(ctx, req)

	return
}

func (v *ValidatorRPC) AbortWithError(ctx context.Context, args ...any) {
	msg := &ext.MiddlewareMessage{
		StatusCode: 500,
		Code:       ext.InternalServerError,
		Desc:       "",
		Err:        nil,
	}

	switch len(args) {
	case 1: // 简化调用：AbortWithError(c, err)
		msg.Code = ext.MissingParameters
		msg.Desc = "元数据错误"
		msg.Err = args[0].(error)
	case 4: // 复杂调用：AbortWithError(c, code, desc, err)
		msg.Code = args[0].(ext.Code)
		msg.Desc = args[1].(string)
		msg.Err = args[2].(error)
	default:
		v.Logger.Error("invalid arguments for AbortWithError")
	}
	// 整理打印
	formattedErr := ext.FormatError(msg.Code, msg.Desc, msg.Err)
	v.Logger.Error(formattedErr)
}
