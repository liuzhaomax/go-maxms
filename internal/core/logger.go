package core

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

// 初始化系统日志
func InitLogger() func() {
	log := GetConfig().Lib.Log
	// TODO 根据时间创建不同的日志文件，减小IO开支
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

// TODO gin日志存入文件 中间件
