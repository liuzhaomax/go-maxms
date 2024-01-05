package business

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/liuzhaomax/go-maxms/src/api_user/model"
	"github.com/liuzhaomax/go-maxms/src/api_user/schema"
	"github.com/redis/go-redis/v9"
	"time"
)

var BusinessUserSet = wire.NewSet(wire.Struct(new(BusinessUser), "*"))

type BusinessUser struct {
	Model *model.ModelUser
	Tx    *core.Trans
	Redis *redis.Client
}

func (b *BusinessUser) PostLogin(c *gin.Context) (string, error) {
	loginReq := &schema.LoginReq{}
	err := c.ShouldBind(loginReq)
	if err != nil {
		return core.EmptyString, core.FormatError(core.ParseIssue, "请求体无效", err)
	}
	user := &model.User{}
	err = b.Tx.ExecTrans(c, func(ctx context.Context) error {
		err = b.Model.QueryUserByUsername(c, loginReq.Username, user)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return core.EmptyString, core.FormatError(core.PermissionDenied, "登录失败", err)
	}
	// loginReq.Password是从SGW经过RSA解码后得到密码
	result := core.VerifyEncodedPwd(loginReq.Password, core.GetConfig().App.Salt, user.Password)
	if !result {
		return core.EmptyString, core.FormatError(core.PermissionDenied, "登录验证失败", err)
	}
	// 生成Bearer jwt，使用userID与ip签发
	j := core.NewJWT()
	duration := time.Second * 60 * 60 * 24 * 7 // 一周
	token, err := j.GenerateToken(user.UserID, core.GetClientIP(c), duration)
	if err != nil {
		return core.EmptyString, core.FormatError(core.PermissionDenied, "Token生成失败", err)
	}
	// 将userID设置到cookie中
	maxAge := int(duration)
	domain := core.GetConfig().App.Domain
	c.SetCookie(
		core.UserID,
		user.UserID,
		maxAge,
		"/",
		domain,
		true,
		true)
	return core.Bearer + token, nil
}

func (b *BusinessUser) DeleteLogin(c *gin.Context) error {
	maxAge := int(time.Millisecond)
	domain := core.GetConfig().App.Domain
	c.SetCookie(
		core.UserID,
		core.EmptyString,
		maxAge,
		"/",
		domain,
		true,
		true)
	return nil
}

func (b *BusinessUser) GetUserByUserID(c *gin.Context) (*schema.UserRes, error) {
	userID := c.Param(core.UserID)
	user := &model.User{}
	err := b.Tx.ExecTrans(c, func(ctx context.Context) error {
		err := b.Model.QueryUserByUserID(c, userID, user)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	userRes := schema.MapUser2UserRes(user)
	return userRes, nil
}
