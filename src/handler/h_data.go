package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var DataSet = wire.NewSet(wire.Struct(new(HData), "*"))

type HData struct {
	//BData *service.BData
	//IRes  core.IResponse
}

func (hData *HData) GetDataById(c *gin.Context) {

}
