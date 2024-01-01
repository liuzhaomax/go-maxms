package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"net/http"
	"os"
	"runtime"
	"time"
)

// 初始化logrus，让在初始化日志前的log也在console中打印，但不记录在日志文件中
func init() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(selectFormatter("text"))
}

// 初始化系统日志
func InitLogger() func() {
	log := GetConfig().Lib.Log
	// TODO NOT NOW 根据时间创建不同的日志文件，减小IO开支
	file, err := os.OpenFile(log.FileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		logrus.WithField(FAILURE, GetFuncName()).Panic(FormatError(IOFailure, "日志文件打开失败", err))
	}
	logger := logrus.New()
	logger.SetLevel(selectLogLevel())
	logger.SetFormatter(selectFormatter("text"))
	//logger.SetReportCaller(true) // 输出caller
	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   log.FileName,
		MaxSize:    50, // megabytes
		MaxBackups: 3,  // amouts
		MaxAge:     28, // days
		Level:      selectLogLevel(),
		Formatter:  selectFormatter(),
	})
	if err != nil {
		logger.WithField(FAILURE, GetFuncName()).Panic(FormatError(Unknown, "日志hook生成失败", err))
		panic(err)
	}
	logger.AddHook(rotateFileHook)
	cfg.App.Logger = logger
	return func() {
		if file != nil {
			err = file.Close()
			if err != nil {
				logger.WithField(FAILURE, GetFuncName()).Panic(FormatError(IOFailure, "日志文件关闭失败", err))
				panic(err)
			}
		}
	}
}

func selectFormatter(forceFormatter ...string) logrus.Formatter {
	log := cfg.Lib.Log
	format := log.Format
	if len(forceFormatter) != 0 {
		format = forceFormatter[0]
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

func LoggerToFile() gin.HandlerFunc {
	logger := cfg.App.Logger
	return func(c *gin.Context) {
		clientIP := GetClientIP(c)
		userAgent := GetUserAgent(c)
		err := ValidateHeaders(c)
		if err != nil {
			cfg.App.Logger.WithField(FAILURE, GetFuncName()).Info(FormatError(MissingParameters, "请求头错误", err))
			c.AbortWithStatusJSON(http.StatusBadRequest, FormatError(MissingParameters, "请求头错误", err))
		}
		LoggerFormat := logrus.Fields{
			"method":     c.Request.Method,
			"uri":        c.Request.RequestURI,
			"client_ip":  clientIP,
			"user_agent": userAgent,
			"token":      c.Request.Header.Get(Authorization),
			"trace_id":   c.Request.Header.Get(TraceId),
			"span_id":    c.Request.Header.Get(SpanId),
			"parent_id":  c.Request.Header.Get(ParentId),
			"app_id":     c.Request.Header.Get(AppId),
		}
		// Incoming日志是来的什么就是什么，只有traceID应一致
		logger.WithFields(LoggerFormat).Info(FormatInfo("请求开始"))
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		took := endTime.Sub(startTime)
		statusCode := c.Writer.Status()

		// json标准写法
		logger.WithFields(LoggerFormat).WithFields(logrus.Fields{
			"took":   took,
			"status": statusCode,
		}).Info(FormatInfo("请求结束"))

		// concatenated json 写法
		//format := &LoggerFormat{
		//    StatusCode: statusCode,
		//    Took:       took,
		//    Method:     c.Request.Method,
		//    URI:        c.Request.RequestURI,
		//    ClientIP:   clientIP,
		//    UserAgent:  userAgent,
		//    Token:      c.Request.Header.Get(Authorization),
		//    TraceId:    c.Request.Header.Get(TraceId),
		//    SpanID:     c.Request.Header.Get(SpanId),
		//    ParentID:   c.Request.Header.Get(ParentId),
		//    AppID:      c.Request.Header.Get(AppId),
		//}
		//formatBytes, _ := json.Marshal(format)
		//logger.Info(string(formatBytes))

		// 竖线分割写法
		//logger.Infof("| %3d | %13v | %15s | %8s | %s | %20s",
		//    statusCode,
		//    took,
		//    clientIP,
		//    method,
		//    uri,
		//    userAgent,
		//)
	}
}

//type LoggerFormat struct {
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
//}
