package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/liuzhaomax/go-maxms/internal/api"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/liuzhaomax/go-maxms/internal/core/config"
	"github.com/liuzhaomax/go-maxms/internal/core/ext"
	"github.com/sirupsen/logrus"
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

func InitConfig(opts *options) {
	cfg := core.GetConfig()
	cfg.LoadConfig(opts.ConfigFile)
	cfg.App.Logger.WithField("path", opts.ConfigFile).Info(ext.FormatInfo("配置文件加载成功"))
	cfg.LoadSecret()
	cfg.App.Logger.WithField("path", ".env").Info(ext.FormatInfo("密钥文件加载成功"))
	cfg.App.Logger.Info(ext.FormatInfo("系统启动"))

	if cfg.App.Enabled.ServiceDiscovery {
		// register service
		err := cfg.Lib.Consul.ServiceRegister()
		if err != nil {
			cfg.App.Logger.WithField(config.FAILURE, ext.GetFuncName()).
				Fatal(ext.FormatError(ext.Unknown, "服务注册失败", err))
		}

		cfg.App.Logger.Info(ext.FormatInfo("服务注册成功"))
		// discover services
		go func() {
			for {
				err = cfg.Lib.Consul.ServiceDiscover()
				if err != nil {
					cfg.App.Logger.WithField(config.FAILURE, ext.GetFuncName()).
						Warn(ext.FormatError(ext.Unknown, "下游服务发现失败", err))
				}

				time.Sleep(time.Duration(cfg.Lib.Consul.Interval) * time.Second)
			}
		}()
	}
}

func InitHttpServer(ctx context.Context, handler http.Handler) func() {
	cfg := core.GetConfig()
	cfg.App.Logger.Info(ext.FormatInfo("服务启动开始"))
	addr := fmt.Sprintf("%s:%s", "0.0.0.0", cfg.Server.Http.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  time.Duration(cfg.Server.Http.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.Http.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.Http.IdleTimeout) * time.Second,
	}

	go func() {
		cfg.App.Logger.WithFields(logrus.Fields{
			"app_name": cfg.App.Name,
			"version":  cfg.App.Version,
			"pid":      os.Getpid(),
			"host":     cfg.Server.Ws.Host,
			"port":     cfg.Server.Ws.Port,
			"protocol": cfg.Server.Ws.Protocol,
		}).Info(ext.FormatInfo("服务启动成功"))

		cfg.App.Logger.WithContext(ctx).Infof("Service %s is running at %s via %s", cfg.App.Name, addr, cfg.Server.Http.Protocol)

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			cfg.App.Logger.WithField(config.FAILURE, ext.GetFuncName()).
				Fatal(ext.FormatError(ext.Unknown, "服务启动失败", err))
		}
	}()

	return func() {
		cfg.App.Logger.Info(ext.FormatInfo("服务关闭开始"))

		_ctx, cancel := context.WithTimeout(
			ctx,
			time.Second*time.Duration(cfg.Server.Http.ShutdownTimeout),
		)
		defer cancel()

		server.SetKeepAlivesEnabled(false)

		err := server.Shutdown(_ctx)
		if err != nil {
			cfg.App.Logger.WithContext(_ctx).
				WithField(config.FAILURE, ext.GetFuncName()).
				Error(ext.FormatError(ext.Unknown, "服务关闭异常", err))
		}

		cfg.App.Logger.Info(ext.FormatInfo("服务关闭成功"))
	}
}

func InitWsServer(ctx context.Context, handler http.Handler) func() {
	cfg := core.GetConfig()
	cfg.App.Logger.Info(ext.FormatInfo("服务启动开始"))
	addr := fmt.Sprintf("%s:%s", "0.0.0.0", cfg.Server.Ws.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  time.Duration(cfg.Server.Ws.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.Ws.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.Ws.IdleTimeout) * time.Second,
	}

	go func() {
		cfg.App.Logger.WithFields(logrus.Fields{
			"app_name": cfg.App.Name,
			"version":  cfg.App.Version,
			"pid":      os.Getpid(),
			"host":     cfg.Server.Ws.Host,
			"port":     cfg.Server.Ws.Port,
			"protocol": cfg.Server.Ws.Protocol,
		}).Info(ext.FormatInfo("服务启动成功"))

		cfg.App.Logger.WithContext(ctx).Infof("Service %s is running at %s via %s", cfg.App.Name, addr, cfg.Server.Ws.Protocol)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			cfg.App.Logger.WithField(config.FAILURE, ext.GetFuncName()).
				Fatal(ext.FormatError(ext.Unknown, "服务启动失败", err))
		}
	}()

	return func() {
		cfg.App.Logger.Info(ext.FormatInfo("服务关闭开始"))

		_ctx, cancel := context.WithTimeout(
			ctx,
			time.Second*time.Duration(cfg.Server.Ws.ShutdownTimeout),
		)
		defer cancel()

		server.SetKeepAlivesEnabled(false)

		if err := server.Shutdown(_ctx); err != nil {
			cfg.App.Logger.WithContext(_ctx).
				WithField(config.FAILURE, ext.GetFuncName()).
				Error(ext.FormatError(ext.Unknown, "服务关闭异常", err))
		}

		cfg.App.Logger.Info(ext.FormatInfo("服务关闭成功"))
	}
}

func InitRpcServer(ctx context.Context, handlerRPC *api.HandlerRPC) func() {
	cfg := core.GetConfig()
	cfg.App.Logger.Info(ext.FormatInfo("服务启动开始"))
	addr := fmt.Sprintf("%s:%s", "0.0.0.0", cfg.Server.Rpc.Port)
	server := handlerRPC.Register()

	go func() {
		listen, err := net.Listen("tcp", addr)
		if err != nil {
			cfg.App.Logger.WithField(config.FAILURE, ext.GetFuncName()).
				Fatal(ext.FormatError(ext.Unknown, "服务监听失败", err))
		}

		cfg.App.Logger.WithContext(ctx).Infof("Service %s is running at %s via %s", cfg.App.Name, addr, cfg.Server.Rpc.Protocol)

		err = server.Serve(listen)
		if err != nil {
			cfg.App.Logger.WithField(config.FAILURE, ext.GetFuncName()).
				Fatal(ext.FormatError(ext.Unknown, "服务启动失败", err))
		}
	}()

	return func() {
		cfg.App.Logger.Info(ext.FormatInfo("服务关闭开始"))

		_, cancel := context.WithTimeout(ctx, time.Second*time.Duration(cfg.Server.Rpc.ShutdownTimeout))
		defer cancel()

		server.Stop()
		cfg.App.Logger.Info(ext.FormatInfo("服务关闭成功"))
	}
}

func Init(ctx context.Context, optFuncs ...Option) func() {
	cfg := core.GetConfig()
	// initialising options
	opts := options{}
	for _, optFunc := range optFuncs {
		optFunc(&opts)
	}
	// init conf
	InitConfig(&opts)

	// init injector
	injector, cleanInjection, _ := InitInjector()
	// init server by protocol
	var cleanServer func()

	switch cfg.Server.Http.Protocol {
	case "http":
		// register apis
		injector.Handler.RegisterStaticFS(injector.InjectorHTTP.Engine, opts.WWWDir) // static
		injector.Handler.Register(injector.InjectorHTTP.Engine)                      // dynamic
		// init server
		cleanServer = InitHttpServer(ctx, injector.InjectorHTTP.Engine)
	case "rpc":
		cleanServer = InitRpcServer(ctx, injector.InjectorRPC.HandlerRPC)
	case "ws":
		injector.HandlerWs.Register(injector.InjectorWS.Engine)
		cleanServer = InitWsServer(ctx, injector.InjectorWS.Engine)
	default:
		cleanServer = InitRpcServer(ctx, injector.HandlerRPC)
	}

	return func() {
		cleanServer()
		cleanInjection()
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
		cfg.App.Logger.WithContext(ctx).Infof("%s [%s]", ext.FormatInfo("服务中断信号收到"), sig.String())

		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			state = 0

			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}

	defer os.Exit(state)
	defer cfg.App.Logger.WithContext(ctx).Infof("%s", ext.FormatInfo("系统停止"))
	defer time.Sleep(time.Second)
	defer clean()
}
