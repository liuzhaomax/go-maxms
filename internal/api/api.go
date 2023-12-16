package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms-template-me/internal/core"
	"github.com/liuzhaomax/go-maxms-template-me/internal/middleware"
	"github.com/liuzhaomax/go-maxms-template-me/src/handler"
	"github.com/liuzhaomax/go-maxms-template-me/src/router"
	"github.com/sirupsen/logrus"
	"net/http"
)

var APISet = wire.NewSet(wire.Struct(new(Handler), "*"), wire.Bind(new(API), new(*Handler)))

type API interface {
	Register(app *gin.Engine)
}

type Handler struct {
	HandlerData *handler.HData
}

func (handler *Handler) Register(app *gin.Engine) {
	app.NoRoute(handler.GetNoRoute)
	app.Use(middleware.Cors())
	app.StaticFS("/static", http.Dir("./static"))
	router.Register(handler.HandlerData, app)
}

func (handler *Handler) GetNoRoute(c *gin.Context) {
	// TODO 封装日志，级别分类（debug有执行方法）
	logrus.WithFields(logrus.Fields{
		"成功方法": core.GetFuncName(),
	}).Info("123")
	c.JSON(http.StatusNotFound, gin.H{"res": "404"})
}
