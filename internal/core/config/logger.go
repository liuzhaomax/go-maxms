package config

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liuzhaomax/go-maxms/internal/core/ext"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// 初始化logrus，让在初始化日志前的log也在console中打印，但不记录在日志文件中
func init() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(selectFormatter("text"))
}

type logConfig struct {
	Level        string `mapstructure:"level"`
	Format       string `mapstructure:"format"`
	Color        bool   `mapstructure:"color"`
	ReportCaller bool   `mapstructure:"report_caller"`
	FilePath     string `mapstructure:"file_path"`
	FileName     string `mapstructure:"file_name"`
}

// 日志扩展loggerx的Provider
func InitLogrus() *logrus.Logger {
	return cfg.App.Logger
}

func InitLogrusEntry() *logrus.Entry {
	return cfg.App.Logger.WithFields(logrus.Fields{})
}

// 初始化系统日志
func InitLogger() *logrus.Logger {
	log := GetConfig().Lib.Log
	if err := os.MkdirAll(log.FilePath, 0o666); err != nil {
		logrus.WithField(FAILURE, ext.GetFuncName()).
			Panic(ext.FormatError(ext.IOException, "日志目录创建失败", err))
	}

	fileName := filepath.Join(log.FilePath, log.FileName)
	logger := logrus.New()
	logger.SetLevel(selectLogLevel())
	logger.SetFormatter(selectFormatter(log.Format))
	logger.SetReportCaller(log.ReportCaller)

	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   fileName,
		MaxSize:    2,               // megabytes
		MaxBackups: 999999999999999, // amounts
		MaxAge:     200,             // days 国家规定日志必须保存6个月，这里设置200天
		Compress:   false,
		Level:      selectLogLevel(),
		Formatter:  selectFormatter(),
	})
	if err != nil {
		logger.WithField(FAILURE, ext.GetFuncName()).
			Panic(ext.FormatError(ext.Unknown, "日志hook生成失败", err))
		panic(err)
	}

	logger.AddHook(rotateFileHook)
	cfg.App.Logger = logger

	return logger
}

func selectFormatter(forceFormatter ...string) logrus.Formatter {
	log := cfg.Lib.Log

	format := log.Format
	if len(forceFormatter) != 0 {
		format = forceFormatter[0]
	}

	if forceFormatter == nil {
		format = "json"
	}

	switch format {
	case "text":
		return &logrus.TextFormatter{
			PadLevelText:    true,
			ForceColors:     log.Color,
			FullTimestamp:   true,
			TimestampFormat: time.DateTime,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				file := fmt.Sprintf("%s:%d", f.File, f.Line)

				return "", fmt.Sprintf("\033[1;34m%s\033[0m", file)
			},
		}
	case "json":
		return &logrus.JSONFormatter{
			TimestampFormat:   time.RFC3339Nano,
			DisableTimestamp:  false,
			DisableHTMLEscape: false,
			DataKey:           "",
			FieldMap:          nil,
			CallerPrettyfier:  nil,
			PrettyPrint:       false,
		}
	default:
		return &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		}
	}
}

func selectLogLevel() logrus.Level {
	switch cfg.Lib.Log.Level {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.InfoLevel
	}
}

func GenGinLoggerFields(c *gin.Context) logrus.Fields {
	return logrus.Fields{
		"method":     c.Request.Method,
		"uri":        c.Request.RequestURI,
		"client_ip":  GetClientIP(c),
		"user_agent": GetUserAgent(c),
		"token":      c.GetHeader(Authorization),
		"trace_id":   c.GetHeader(TraceId),
		"span_id":    c.GetHeader(SpanId),
		"parent_id":  c.GetHeader(ParentId),
		"app_id":     c.GetHeader(AppId),
		"request_id": c.GetHeader(RequestId),
		"user_id":    c.GetHeader(UserId),
	}
}

// HTTP 日志中间件
func LoggerForHTTP() gin.HandlerFunc {
	logger := cfg.App.Logger

	return func(c *gin.Context) {
		// 过滤ws心跳
		// if shouldSkipHeartbeatLogging(c) {
		// 	c.Next()
		// 	return
		// }
		loggerFormat := GenGinLoggerFields(c)
		// Incoming日志是来的什么就是什么，只有traceID应一致
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		took := endTime.Sub(startTime).Milliseconds()
		statusCode := c.Writer.Status()
		// json标准写法
		logger.WithFields(loggerFormat).WithFields(logrus.Fields{
			"took":   took,
			"status": statusCode,
		}).Info("请求结束")
		// concatenated json 写法
		//
		//	format := &LoggerFormat{
		//	   StatusCode: statusCode,
		//	   Took:       took,
		//	   Method:     c.Request.Method,
		//	   URI:        c.Request.RequestURI,
		//	   ClientIP:   clientIP,
		//	   UserAgent:  userAgent,
		//	   Token:      c.Request.Header.Get(Authorization),
		//	   TraceId:    c.Request.Header.Get(TraceId),
		//	   SpanID:     c.Request.Header.Get(SpanId),
		//	   ParentID:   c.Request.Header.Get(ParentId),
		//	   AppID:      c.Request.Header.Get(AppId),
		//	   RequestID:  c.Request.Header.Get(RequestId),
		//	   UserID:     c.Request.Header.Get(UserId),
		//	}
		//
		// formatBytes, _ := json.Marshal(format)
		// logger.Info(string(formatBytes))
		// 竖线分割写法
		// logger.Infof("| %3d | %13v | %15s | %8s | %s | %20s",
		//
		//	statusCode,
		//	took,
		//	clientIP,
		//	method,
		//	uri,
		//	userAgent,
		//
		// )
	}
}

// type LoggerFormat struct {
//    StatusCode int           `json:"code"`
//    Took       time.Duration `json:"took"`
//    Method     string        `json:"method"`
//    URI        string        `json:"uri"`
//    ClientIP   string        `json:"client_ip"`
//    UserAgent  string        `json:"user_agent"`
//    Token      string        `json:"token"`
//    TraceId    string        `json:"trace_id"`
//    SpanID     string        `json:"span_id"`
//    ParentID   string        `json:"parent_id"`
//    UpstreamID string        `json:"upstream_id"`
//    AppID      string        `json:"app_id"`
//    RequestID  string        `json:"request_id"`
//    UserID  string        `json:"user_id"`
// }

// RPC 日志中间件
func LoggerForRPC(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	logger := cfg.App.Logger

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		LogFailure(ext.NotFound, "缺少metadata", nil)
	}

	LoggerFormat := logrus.Fields{
		"method":     SelectFromMetadata(md, Method),
		"uri":        SelectFromMetadata(md, RequestURI),
		"client_ip":  SelectFromMetadata(md, ClientIp),
		"user_agent": SelectFromMetadata(md, UserAgent),
		"token":      SelectFromMetadata(md, Authorization),
		"trace_id":   SelectFromMetadata(md, TraceId),
		"span_id":    SelectFromMetadata(md, SpanId),
		"parent_id":  SelectFromMetadata(md, ParentId),
		"app_id":     SelectFromMetadata(md, AppId),
		"request_id": SelectFromMetadata(md, RequestId),
		"user_id":    SelectFromMetadata(md, UserId),
	}
	logger.WithFields(LoggerFormat).Info("请求开始")

	startTime := time.Now()
	res, err := handler(ctx, req)
	endTime := time.Now()
	took := endTime.Sub(startTime).Milliseconds()
	// json标准写法
	logger.WithFields(LoggerFormat).WithFields(logrus.Fields{
		"took": took,
	}).Info("请求结束")

	return res, err
}

func LogSuccess(desc string) {
	cfg.App.Logger.Info(ext.FormatInfo(desc))
}

func LogFailure(code ext.Code, desc string, err error) {
	cfg.App.Logger.Error(ext.FormatError(code, desc, err))
}

// shouldSkipHeartbeatLogging 过滤ws心跳
func shouldSkipHeartbeatLogging(c *gin.Context) bool {
	if !strings.Contains(c.Request.RequestURI, "/ws") {
		return false
	}
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return false
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var msg struct {
		Action string `json:"action"`
		Body   string `json:"body"`
	}
	if json.Unmarshal(bodyBytes, &msg) != nil {
		return false
	}

	return msg.Action == "ping"
}
