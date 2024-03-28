package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/liuzhaomax/go-maxms/src/api_user/business"
)

var HandlerUserSet = wire.NewSet(wire.Struct(new(HandlerUser), "*"))

type HandlerUser struct {
	Business    *business.BusinessUser
	Logger      core.ILogger
	Res         core.IResponse
	RocketMQSet core.IRocketMQ
}

func (h *HandlerUser) GetPuk(c *gin.Context) {
	h.Res.ResSuccess(c, core.GetConfig().App.PublicKeyStr)
}

func (h *HandlerUser) PostLogin(c *gin.Context) {
	tokenRes, err := h.Business.PostLogin(c)
	if err != nil {
		h.Res.ResFailure(c, 500, core.Unknown, "登录失败", err)
		return
	}
	h.Res.ResSuccess(c, tokenRes)
}

func (h *HandlerUser) DeleteLogin(c *gin.Context) {
	err := h.Business.DeleteLogin(c)
	if err != nil {
		h.Res.ResFailure(c, 500, core.Unknown, "登出失败", err)
		return
	}
	h.Res.ResSuccess(c, nil)
}

func (h *HandlerUser) GetUserByUserID(c *gin.Context) {
	user, err := h.Business.GetUserByUserID(c)
	if err != nil {
		h.Res.ResFailure(c, 500, core.Unknown, "查询失败", err)
		return
	}
	h.Res.ResSuccess(c, user)
}
