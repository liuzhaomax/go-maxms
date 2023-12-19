package core

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

// 初始化系统日志
func InitLogger() func() {
	log := GetConfig().Lib.Log
	// TODO NOT NOW 根据时间创建不同的日志文件，减小IO开支
	file, err := os.OpenFile(log.FileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.WithField("失败方法", GetFuncName()).Panic("日志文件创建或打开失败")
		panic(err)
	}
	logrus.SetLevel(selectLogLevel(&log))
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: log.Color})
	logrus.SetOutput(io.MultiWriter(file, os.Stdout))
	return func() {
		if file != nil {
			err = file.Close()
			if err != nil {
				logrus.WithField("失败方法", GetFuncName()).Panic("日志文件关闭失败")
				panic(err)
			}
		}
	}
}

func selectLogLevel(log *Log) logrus.Level {
	switch log.Level {
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
}

func LoggerToFile() (gin.HandlerFunc, *logrus.Logger) {
	log := GetConfig().Lib.Log
	src, err := os.OpenFile(log.FileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		logrus.WithField("失败方法", GetFuncName()).Panic(FormatError(Unknown, "日志文件打开失败", err))
	}
	logger := logrus.New()
	logger.SetLevel(selectLogLevel(&log))
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   time.RFC3339Nano,
		DisableTimestamp:  false,
		DisableHTMLEscape: false,
		DataKey:           "",
		FieldMap:          nil,
		CallerPrettyfier:  nil,
		PrettyPrint:       false,
	})
	logger.Out = src
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		took := endTime.Sub(startTime)
		method := c.Request.Method
		uri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		// 竖线分割写法
		//logger.Infof("| %3d | %13v | %15s | %8s | %s ",
		//    statusCode,
		//    took,
		//    clientIP,
		//    method,
		//    uri,
		//)

		// concatenated json 写法
		format := &LoggerFormat{
			StatusCode: statusCode,
			Took:       took,
			ClientIP:   clientIP,
			Method:     method,
			URI:        uri,
		}
		formatBytes, _ := json.Marshal(format)
		logger.Info(string(formatBytes))

		// json标准写法
		//logger.WithFields(logrus.Fields{
		//	"code":      statusCode,
		//	"took":      took,
		//	"client_ip": clientIP,
		//	"method":    method,
		//	"uri":       uri,
		//}).Info("123")
	}, logger
}
