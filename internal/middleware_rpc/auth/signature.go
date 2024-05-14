package auth

import (
	"context"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

func (auth *AuthRPC) ValidateSignature(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	cfg := core.GetConfig()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = auth.GenErrMsg(ctx, "元数据解析错误", err)
		return
	}
	if len(md[core.UserId]) == 0 || len(md[core.SpanId]) == 0 || len(md[core.RequestURI]) == 0 {
		err = auth.GenErrMsg(ctx, "签名信息缺失", err)
		return
	}
	userId := md[core.UserId][0]
	nonceForValidation := md[core.ParentId][0]
	nonce := md[core.SpanId][0]
	if nonce == core.EmptyString {
		nonce = core.ShortUUID()
	}
	// 根据headers里给定的信息，生成签名并比对
	signatureGen := core.GenAppSignature(cfg.App.Id, cfg.App.Secret, userId, nonceForValidation)
	signatureMD := md[core.Signature][0]
	if signatureGen != signatureMD {
		err = auth.GenErrMsg(ctx, "签名验证失败", errors.New("签名不匹配"))
		return
	}
	result, err := auth.Redis.SAdd(context.Background(), core.Nonce, nonce).Result()
	// 如果直接使用返回值，(*result).Val()，1是set里原来没有，加入成功，0是set里原来有，加入失败
	if err != nil {
		err = auth.GenErrMsg(ctx, "签名验证失败", err)
		return
	}
	if result == 0 {
		err = auth.GenErrMsg(ctx, "签名验证失败", errors.New("set已有该值"))
		return
	}
	// 设置过期时间
	err = auth.Redis.Expire(context.Background(), core.Nonce, time.Second*5).Err()
	if err != nil {
		err = auth.GenErrMsg(ctx, "签名过期时间设置失败", err)
		return
	}
	auth.GenOkMsg(ctx, "签名已写入缓存")
	resp, err = handler(ctx, req)
	return
}
