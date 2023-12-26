package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/liuzhaomax/go-maxms/src/data_api/business"
	"github.com/sirupsen/logrus"
)

var HandlerDataSet = wire.NewSet(wire.Struct(new(HandlerData), "*"))

type HandlerData struct {
	Business *business.BusinessData
	Logger   *logrus.Logger
	IRes     core.IResponse
}

func (h *HandlerData) GetDataById(c *gin.Context) {
	h.IRes.ResSuccess(c, core.GetFuncName(), "ok")
	//c.JSON(http.StatusOK, gin.H{"res": "ok"})
}
