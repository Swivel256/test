package main

import (
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

func main() {

	// 创建 logrus 实例
	log := logrus.New()

	// 设置日志级别
	log.SetLevel(logrus.DebugLevel)


	log.SetFormatter(&logrus.JSONFormatter{})
	logDir := "log"
	if err := os.MkdirAll(logDir, 0777); err != nil {
		log.Panic(err)
	}
	writer, err := rotatelogs.New(
		path.Join(logDir, "access_log.%Y%m%d"),
		rotatelogs.WithLinkName("access_log"),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	if err != nil {
		log.Panic(err)
	}

	log.AddHook(
		lfshook.NewHook(
			lfshook.WriterMap{
				logrus.DebugLevel: writer,
				logrus.InfoLevel:  writer,
				logrus.WarnLevel:  writer,
				logrus.ErrorLevel: writer,
				logrus.FatalLevel: writer,
			}, &logrus.JSONFormatter{}),
	)

	// 输出info级别日志
	log.Info("info log")

	// 输出warn级别日志
	log.Warn("warn log")
}
