package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms-template/internal/middleware"
	"github.com/liuzhaomax/go-maxms-template/internal/middleware/cors"
	"github.com/liuzhaomax/go-maxms-template/src/data_api/handler"
	"github.com/liuzhaomax/go-maxms-template/src/router"
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

func (h *Handler) Register(app *gin.Engine) {
	app.NoRoute(h.GetNoRoute)
	app.Use(cors.Cors())
	router.Register(app, h.HandlerData, h.Middleware)
}

func (h *Handler) RegisterStaticFS(app *gin.Engine, path string) {
	app.StaticFS("/"+path, http.Dir("./"+path))
}

func (h *Handler) GetNoRoute(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"res": "404"})
}
