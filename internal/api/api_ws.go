package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/liuzhaomax/go-maxms/internal/core/config"
	"github.com/liuzhaomax/go-maxms/internal/middleware"
	"github.com/liuzhaomax/go-maxms/internal/middleware/cors"
	"github.com/liuzhaomax/go-maxms/src/api_user/handler"
	"github.com/liuzhaomax/go-maxms/src/api_user/router"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var APIWSSet = wire.NewSet(wire.Struct(new(HandlerWs), "*"), wire.Bind(new(API), new(*HandlerWs)))

type HandlerWs struct {
	Middleware         *middleware.Middleware
	Handler            *handler.HandlerUser
	PrometheusRegistry *prometheus.Registry
}

func (h *HandlerWs) Register(app *gin.Engine) {
	cfg := core.GetConfig()
	// 404
	app.NoRoute(h.GetNoRoute)
	// CORS
	app.Use(cors.Cors())
	// consul
	if cfg.App.Enabled.HealthCheck {
		app.GET("/health", h.HealthHandler)
	}
	// prometheus
	if cfg.App.Enabled.Prometheus {
		app.GET("/metrics", h.MetricsHandler)
	}
	// jaeger
	if cfg.App.Enabled.Jaeger {
		app.Use(h.Middleware.Tracing.Trace())
	}
	// 日志
	app.Use(config.LoggerForHTTP())
	// root route
	root := app.Group(cfg.Server.Ws.BaseUrl)
	{
		// interceptor
		if cfg.App.Enabled.HeaderParams {
			root.Use(h.Middleware.Validator.ValidateHeaders())
		}

		if cfg.App.Enabled.Signature {
			root.Use(h.Middleware.Auth.ValidateSignature())
		}
		// dynamic api
		router.RegisterWs(root, h.Handler, h.Middleware)
	}
}

func (h *HandlerWs) RegisterStaticFS(app *gin.Engine, path string) {
	app.StaticFS("/"+path, http.Dir("./"+path))
}

func (h *HandlerWs) GetNoRoute(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"res": "404"})
}

func (h *HandlerWs) HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"health": "ok"})
}

func (h *HandlerWs) MetricsHandler(c *gin.Context) {
	promhttp.HandlerFor(h.PrometheusRegistry, promhttp.HandlerOpts{
		Registry: h.PrometheusRegistry,
	}).ServeHTTP(c.Writer, c.Request)
}
