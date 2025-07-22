package api

import (
	"github.com/liuzhaomax/go-maxms/internal/core/ext"
	"github.com/sirupsen/logrus"
	"net/http"

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
)

var APISet = wire.NewSet(wire.Struct(new(Handler), "*"), wire.Bind(new(API), new(*Handler)))

type API interface {
	Register(app *gin.Engine)
}

type Handler struct {
	Middleware         *middleware.Middleware
	HandlerUser        *handler.HandlerUser
	PrometheusRegistry *prometheus.Registry
}

func (h *Handler) Register(app *gin.Engine) {
	cfg := core.GetConfig()
	// 404
	app.NoRoute(h.GetNoRoute)
	// root route
	root := app.Group(cfg.Server.Http.BaseUrl)
	{
		// CORS
		root.Use(cors.Cors())
		// consul
		if cfg.App.Enabled.HealthCheck {
			root.GET("/health", h.HealthHandler)
		}
		// prometheus
		if cfg.App.Enabled.Prometheus {
			root.GET("/metrics", h.MetricsHandler)
		}
		// jaeger
		if cfg.App.Enabled.Jaeger {
			root.Use(h.Middleware.Tracing.Trace())
		}
		// 日志
		root.Use(config.LoggerForHTTP())
		// interceptor
		if cfg.App.Enabled.HeaderParams {
			root.Use(h.Middleware.Validator.ValidateHeaders())
		}

		if cfg.App.Enabled.Signature {
			root.Use(h.Middleware.Auth.ValidateSignature())
		}
		// dynamic api
		router.Register(root, h.HandlerUser, h.Middleware)
	}
}

func (h *Handler) RegisterStaticFS(app *gin.Engine, path string) {
	app.StaticFS("/"+path, http.Dir("./"+path))
}

func (h *Handler) GetNoRoute(c *gin.Context) {
	cfg := core.GetConfig()
	LoggerFormat := logrus.Fields{
		"method":     c.Request.Method,
		"uri":        c.Request.RequestURI,
		"client_ip":  config.GetClientIP(c),
		"user_agent": config.GetUserAgent(c),
		"token":      c.GetHeader(config.Authorization),
		"trace_id":   c.GetHeader(config.TraceId),
		"span_id":    c.GetHeader(config.SpanId),
		"parent_id":  c.GetHeader(config.ParentId),
		"app_id":     c.GetHeader(config.AppId),
		"request_id": c.GetHeader(config.RequestId),
		"user_id":    c.GetHeader(config.UserId),
	}
	cfg.App.Logger.WithFields(LoggerFormat).Error(ext.FormatError(ext.NotFound, "错误路径", nil))
	c.JSON(http.StatusNotFound, gin.H{"res": "404"})
}

func (h *Handler) HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"health": "ok"})
}

func (h *Handler) MetricsHandler(c *gin.Context) {
	promhttp.HandlerFor(h.PrometheusRegistry, promhttp.HandlerOpts{
		Registry: h.PrometheusRegistry,
	}).ServeHTTP(c.Writer, c.Request)
}
