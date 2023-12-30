package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/liuzhaomax/go-maxms/src/api_user/business"
	"github.com/sirupsen/logrus"
)

var HandlerUserSet = wire.NewSet(wire.Struct(new(HandlerUser), "*"))

type HandlerUser struct {
	Business *business.BusinessUser
	Logger   *logrus.Logger
	IRes     core.IResponse
}

func (h *HandlerUser) GetPuk(c *gin.Context) {
	h.IRes.ResSuccess(c, core.GetFuncName(), core.GetConfig().App.PublicKeyStr)
}

func (h *HandlerUser) PostLogin(c *gin.Context) {
	token, err := h.Business.PostLogin(c)
	if err != nil {
		h.IRes.ResFailure(c, core.GetFuncName(), 500, core.Unknown, "登录失败", err)
		return
	}
	h.IRes.ResSuccess(c, core.GetFuncName(), token)
}

func (h *HandlerUser) GetUserByUserID(c *gin.Context) {
	user, err := h.Business.GetUserByUserID(c)
	if err != nil {
		h.IRes.ResFailure(c, core.GetFuncName(), 500, core.Unknown, "查询失败", err)
		return
	}
	h.IRes.ResSuccess(c, core.GetFuncName(), user)
}
