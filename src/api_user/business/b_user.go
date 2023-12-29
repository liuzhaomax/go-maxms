package business

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/liuzhaomax/go-maxms/src/api_user/model"
	"github.com/liuzhaomax/go-maxms/src/api_user/schema"
	"github.com/sirupsen/logrus"
)

var BusinessUserSet = wire.NewSet(wire.Struct(new(BusinessUser), "*"))

type BusinessUser struct {
	Model  *model.ModelUser
	Logger *logrus.Logger
	Tx     *core.Trans
}

func (b *BusinessUser) GetUserByUserID(c *gin.Context) (*schema.UserRes, error) {
	userID := c.Param("userID")
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
