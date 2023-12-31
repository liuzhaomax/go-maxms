package business

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/liuzhaomax/go-maxms/src/api_user/model"
	"github.com/liuzhaomax/go-maxms/src/api_user/schema"
	"github.com/liuzhaomax/go-maxms/src/utils"
	"github.com/sirupsen/logrus"
	"time"
)

var BusinessUserSet = wire.NewSet(wire.Struct(new(BusinessUser), "*"))

type BusinessUser struct {
	Model  *model.ModelUser
	Logger *logrus.Logger
	Tx     *core.Trans
}

func (b *BusinessUser) PostLogin(c *gin.Context) (string, error) {
	loginReq := &schema.LoginReq{}
	err := c.ShouldBind(loginReq)
	if err != nil {
		return core.EmptyString, core.FormatError(core.ParseIssue, "请求体无效", err)
	}
	user := &model.User{}
	err = b.Tx.ExecTrans(c, func(ctx context.Context) error {
		err = b.Model.QueryUserByUsername(loginReq.Username, user)
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
	if result == false {
		return core.EmptyString, core.FormatError(core.PermissionDenied, "登录验证失败", err)
	}
	// 生成Bearer jwt，使用userID与ip签发
	j := core.NewJWT()
	duration := time.Second * 60 * 60 * 24 * 7 // 一周
	token, err := j.GenerateToken(user.UserID, c.ClientIP(), duration)
	if err != nil {
		return core.EmptyString, core.FormatError(core.PermissionDenied, "Token生成失败", err)
	}
	// 将userID设置到cookie中
	maxAge := int(duration)
	domain := core.GetConfig().App.Domain
	c.SetCookie(
		"userID",
		user.UserID,
		maxAge,
		"/",
		domain,
		true,
		true)
	return utils.Bearer + token, nil
}

func (b *BusinessUser) GetUserByUserID(c *gin.Context) (*schema.UserRes, error) {
	userID := c.Param(utils.UserID)
	user := &model.User{}
	err := b.Tx.ExecTrans(c, func(ctx context.Context) error {
		err := b.Model.QueryUserByUserID(userID, user)
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
