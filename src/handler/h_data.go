package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms-template-me/internal/core"
)

var DataSet = wire.NewSet(wire.Struct(new(HData), "*"))

type HData struct {
	IRes core.IResponse
	//BData *service.BData
}

func (hData *HData) GetDataById(c *gin.Context) {
	hData.IRes.ResSuccess(c, core.GetFuncName(), "ok")
	//c.JSON(http.StatusOK, gin.H{"res": "ok"})
}
