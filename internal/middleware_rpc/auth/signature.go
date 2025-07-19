package auth

import (
	"context"
	"time"

	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/liuzhaomax/go-maxms/internal/core/config"
	"github.com/liuzhaomax/go-maxms/internal/core/ext"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (auth *AuthRPC) ValidateSignature(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	cfg := core.GetConfig()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		auth.AbortWithError(ctx, errors.New("元数据解析错误"))

		return resp, err
	}

	if len(md[config.UserId]) == 0 || len(md[config.SpanId]) == 0 ||
		len(md[config.RequestURI]) == 0 {
		auth.AbortWithError(ctx, errors.New("签名信息缺失"))

		return resp, err
	}

	userId := md[config.UserId][0]
	nonceForValidation := md[config.ParentId][0]

	nonce := md[config.SpanId][0]
	if nonce == "" {
		nonce = ext.ShortUUID()
	}
	// 根据headers里给定的信息，生成签名并比对
	signatureGen := ext.GenAppSignature(cfg.App.Id, cfg.App.Secret, userId, nonceForValidation)

	signatureMD := md[config.Signature][0]

	if signatureGen != signatureMD {
		auth.AbortWithError(ctx, ext.Unauthorized, "签名验证失败", errors.New("签名不匹配"))

		return resp, err
	}

	result, err := auth.Redis.SAdd(context.Background(), config.Nonce, nonce).Result()
	// 如果直接使用返回值，(*result).Val()，1是set里原来没有，加入成功，0是set里原来有，加入失败
	if err != nil {
		auth.AbortWithError(ctx, ext.Unauthorized, "签名验证失败", err)

		return resp, err
	}

	if result == 0 {
		auth.AbortWithError(ctx, ext.Unauthorized, "签名验证失败", errors.New("set已有该值"))

		return resp, err
	}
	// 设置过期时间
	err = auth.Redis.Expire(context.Background(), config.Nonce, time.Second*5).Err()
	if err != nil {
		auth.AbortWithError(ctx, ext.Unauthorized, "签名过期时间设置失败", err)

		return resp, err
	}

	resp, err = handler(ctx, req)

	return resp, err
}
