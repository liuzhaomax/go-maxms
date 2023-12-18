package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms-template-me/internal/middleware"
	"github.com/liuzhaomax/go-maxms-template-me/internal/middleware/cors"
	"github.com/liuzhaomax/go-maxms-template-me/src/dataAPI/handler"
	"github.com/liuzhaomax/go-maxms-template-me/src/router"
	"net/http"
)

var APISet = wire.NewSet(wire.Struct(new(Handler), "*"), wire.Bind(new(API), new(*Handler)))

type API interface {
	Register(app *gin.Engine)
}

type Handler struct {
	Middleware  *middleware.Middleware
	HandlerData *handler.HandlerData
}

func (handler *Handler) Register(app *gin.Engine) {
	app.NoRoute(handler.GetNoRoute)
	app.Use(cors.Cors())
	app.StaticFS("/static", http.Dir("./static"))
	router.Register(app, handler.HandlerData, handler.Middleware)
}

func (handler *Handler) GetNoRoute(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"res": "404"})
}
