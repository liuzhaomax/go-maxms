package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms-template-me/internal/core"
	"github.com/sirupsen/logrus"
)

var HandlerDataSet = wire.NewSet(wire.Struct(new(HandlerData), "*"))

type HandlerData struct {
	IRes   core.IResponse
	Logger *logrus.Logger
	//BData *service.BData
}

func (h *HandlerData) GetDataById(c *gin.Context) {
	h.IRes.ResSuccess(c, core.GetFuncName(), "ok")
	//c.JSON(http.StatusOK, gin.H{"res": "ok"})
}
