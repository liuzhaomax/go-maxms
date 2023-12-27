package app

import (
	"context"
	"fmt"
	"github.com/liuzhaomax/go-maxms/internal/core"
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
	cfg := core.GetConfig()
	cleanLogger := cfg.LoadConfig(opts.ConfigFile)
	cfg.App.Logger.WithField("path", opts.ConfigFile).Info(core.FormatInfo("配置文件加载成功"))
	return func() {
		cleanLogger()
	}
}

func InitServer(ctx context.Context, handler http.Handler) func() {
	cfg := core.GetConfig()
	cfg.App.Logger.Info(core.FormatInfo("服务启动开始"))
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}
	go func() {
		cfg.App.Logger.WithContext(ctx).Infof("Service is running at %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			cfg.App.Logger.WithField(core.FAILURE, core.GetFuncName()).Fatal(core.FormatError(core.Unknown, "服务启动失败", err))
		}
	}()
	return func() {
		cfg.App.Logger.Info(core.FormatInfo("服务关闭开始"))
		_ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(cfg.Server.ShutdownTimeout))
		defer cancel()
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(_ctx); err != nil {
			cfg.App.Logger.WithContext(_ctx).WithField(core.FAILURE, core.GetFuncName()).Error(core.FormatError(1, "服务关闭异常", err))
		}
		cfg.App.Logger.Info(core.FormatInfo("服务关闭成功"))
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
	injector, cleanInjection, _ := InitInjector()
	// register apis
	injector.Handler.RegisterStaticFS(injector.Engine, opts.WWWDir) // static
	injector.Handler.Register(injector.Engine)                      // dynamic
	// init server
	cleanServer := InitServer(ctx, injector.Engine)
	cfg := core.GetConfig()
	cfg.App.Logger.WithFields(logrus.Fields{
		"app_name": cfg.App.Name,
		"version":  cfg.App.Version,
		"pid":      os.Getpid(),
		"host":     cfg.Server.Host,
		"port":     cfg.Server.Port,
	}).Info(core.FormatInfo("服务启动成功"))
	return func() {
		cleanServer()
		cleanInjection()
		cleanConfig()
	}
}

func Launch(ctx context.Context, opts ...Option) {
	cfg := core.GetConfig()
	clean := Init(ctx, opts...)
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
EXIT:
	for {
		sig := <-sc
		cfg.App.Logger.WithContext(ctx).Infof("%s [%s]", core.FormatInfo("服务中断信号收到"), sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			state = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}
	defer time.Sleep(time.Second)
	defer os.Exit(state)
	defer clean()
	defer cfg.App.Logger.WithContext(ctx).Infof(core.FormatInfo("系统终止运行"))
}
