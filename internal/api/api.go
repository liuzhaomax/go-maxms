package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/middleware"
	"github.com/liuzhaomax/go-maxms/internal/middleware/cors"
	"github.com/liuzhaomax/go-maxms/src/api_user/handler"
	"github.com/liuzhaomax/go-maxms/src/router"
	"net/http"
)

var APISet = wire.NewSet(wire.Struct(new(Handler), "*"), wire.Bind(new(API), new(*Handler)))

type API interface {
	Register(app *gin.Engine)
}

type Handler struct {
	Middleware  *middleware.Middleware
	HandlerUser *handler.HandlerUser
}

func (h *Handler) Register(app *gin.Engine) {
	app.NoRoute(h.GetNoRoute)
	app.Use(cors.Cors())
	app.Use(h.Middleware.Auth.ValidateSignature())
	router.Register(app, h.HandlerUser, h.Middleware)
}

func (h *Handler) RegisterStaticFS(app *gin.Engine, path string) {
	app.StaticFS("/"+path, http.Dir("./"+path))
}

func (h *Handler) GetNoRoute(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"res": "404"})
}
