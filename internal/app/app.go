package app

import (
	"context"
	"fmt"
	"github.com/liuzhaomax/go-maxms-template-me/internal/core"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Option func(*options)

type options struct {
	ConfigFile string
	WWWDir     string
}

func SetConfigFile(configFile string) Option {
	return func(opts *options) {
		opts.ConfigFile = configFile
	}
}

func SetWWWDir(wwwDir string) Option {
	return func(opts *options) {
		opts.WWWDir = wwwDir
	}
}

func InitConfig(opts *options) func() {
	core.GetConfig().LoadConfig(opts.ConfigFile)
	cleanLogger := core.InitLogger()
	logrus.WithField("path", opts.ConfigFile).Info(core.FormatInfo("配置文件加载成功"))
	return func() {
		cleanLogger()
	}
}

func InitServer(ctx context.Context, handler http.Handler) func() {
	logrus.Info(core.FormatInfo("服务启动开始"))
	cfg := core.GetConfig()
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithField("失败方法", core.GetFuncName()).Fatal(core.FormatError(core.Unknown, "服务启动失败", err))
		}
		logrus.WithFields(logrus.Fields{
			"app_name": cfg.App.Name,
			"version":  cfg.App.Version,
			"pid":      os.Getpid(),
			"host":     cfg.Server.Host,
			"port":     cfg.Server.Port,
		}).Info(core.FormatInfo("服务启动成功"))
		logrus.WithContext(ctx).Infof("Server is running at %s", addr)
	}()
	return func() {
		logrus.Info(core.FormatInfo("服务关闭开始"))
		_ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(cfg.Server.ShutdownTimeout))
		defer cancel()
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(_ctx); err != nil {
			logrus.WithContext(_ctx).WithField("失败方法", core.GetFuncName()).Error(core.FormatError(1, "服务关闭异常", err))
		}
		logrus.Info(core.FormatInfo("服务关闭成功"))
	}
}

func Init(ctx context.Context, optFuncs ...Option) func() {
	// initialising options
	opts := options{}
	for _, optFunc := range optFuncs {
		optFunc(&opts)
	}
	// init conf
	cleanConfig := InitConfig(&opts)
	// init injector
	injector, _ := InitInjector()
	// init server
	cleanServer := InitServer(ctx, injector.Engine)
	return func() {
		cleanConfig()
		cleanServer()
	}
}

func Launch(ctx context.Context, opts ...Option) {
	clean := Init(ctx, opts...)
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
EXIT:
	for {
		sig := <-sc
		logrus.WithContext(ctx).Infof("%s [%s]", core.FormatInfo("服务中断信号收到"), sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			state = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}
	defer logrus.WithContext(ctx).Infof(core.FormatInfo("服务正在关闭"))
	defer time.Sleep(time.Second)
	defer os.Exit(state)
	defer clean()
}
