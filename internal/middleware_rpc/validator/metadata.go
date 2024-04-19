package validator

import (
	"context"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var ValidatorRPCSet = wire.NewSet(wire.Struct(new(ValidatorRPC), "*"))

type ValidatorRPC struct {
	Logger core.ILogger
	Redis  *redis.Client
}

func (v *ValidatorRPC) ValidateMetadata(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = v.GenErrMsg(ctx, "元数据解析错误", err)
		return
	}
	err = core.ValidateMetadata(md)
	if err != nil {
		err = v.GenErrMsg(ctx, "元数据错误", err)
		return
	}
	resp, err = handler(ctx, req)
	return
}

func (v *ValidatorRPC) GenOkMsg(ctx context.Context, desc string) string {
	v.Logger.SucceedWithFieldForRPC(ctx, desc)
	return core.FormatInfo(desc)
}

func (v *ValidatorRPC) GenErrMsg(ctx context.Context, desc string, err error) error {
	v.Logger.FailWithFieldForRPC(ctx, core.MissingParameters, desc, err)
	return core.FormatError(core.MissingParameters, desc, err)
}
