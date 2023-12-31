package core

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"net/http"
	"os"
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

type LoggerFormat struct {
	StatusCode int           `json:"code"`
	Took       time.Duration `json:"took"`
	ClientIP   string        `json:"client_ip"`
	Method     string        `json:"method"`
	URI        string        `json:"uri"`
	UserAgent  string        `json:"user_agent"`
	TraceID    string        `json:"trace_id"`
	SpanID     string        `json:"span_id"`
	ParentID   string        `json:"parent_id"`
}

func LoggerToFile() gin.HandlerFunc {
	logger := cfg.App.Logger
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		method := c.Request.Method
		uri := c.Request.RequestURI
		err := SetHeaders(c)
		traceID := c.Request.Header.Get(TraceId)
		spanID := c.Request.Header.Get(SpanId)
		parentID := c.Request.Header.Get(ParentId)
		logger.WithFields(logrus.Fields{
			"client_ip":  clientIP,
			"user-agent": userAgent,
			"uri":        uri,
			"method":     method,
			"trace_id":   traceID,
			"span_id":    spanID,
			"parent_id":  parentID,
		}).Info("Request Incoming")
		if err != nil {
			cfg.App.Logger.WithField(FAILURE, GetFuncName()).Info(FormatError(MissingParameters, "请求头错误", err))
			c.AbortWithStatusJSON(http.StatusBadRequest, FormatError(MissingParameters, "请求头错误", err))
		}
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		took := endTime.Sub(startTime)
		statusCode := c.Writer.Status()

		// concatenated json 写法
		format := &LoggerFormat{
			StatusCode: statusCode,
			Took:       took,
			ClientIP:   clientIP,
			Method:     method,
			URI:        uri,
			UserAgent:  userAgent,
			TraceID:    traceID,
			SpanID:     spanID,
			ParentID:   parentID,
		}
		formatBytes, _ := json.Marshal(format)
		logger.Info(string(formatBytes))

		// 竖线分割写法
		//logger.Infof("| %3d | %13v | %15s | %8s | %s | %20s",
		//    statusCode,
		//    took,
		//    clientIP,
		//    method,
		//    uri,
		//    userAgent,
		//)

		// json标准写法
		//logger.WithFields(logrus.Fields{
		//	"code":       statusCode,
		//	"took":       took,
		//	"client_ip":  clientIP,
		//	"method":     method,
		//	"uri":        uri,
		//  "user_agent": userAgent,
		//}).Info("123")
	}
}

func SetHeaders(c *gin.Context) error {
	if c.Request.Header.Get(TraceId) == EmptyString || c.Request.Header.Get(SpanId) == EmptyString {
		return errors.New("缺失链路信息")
	}
	c.Request.Header.Set(TraceId, c.Request.Header.Get(TraceId))
	c.Request.Header.Set(ParentId, c.Request.Header.Get(SpanId))
	c.Request.Header.Set(SpanId, SpanID())
	c.Request.Header.Set(ClientIp, c.ClientIP())
	return nil
}