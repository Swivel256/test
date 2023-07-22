package main

import (
	"context"
	"fmt"
	"github.com/ginvmbot/aitrade/pkg"
	"github.com/ginvmbot/aitrade/pkg/cmdutil"
	"github.com/ginvmbot/aitrade/pkg/filter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"

	"os"
	"path"
	"syscall"
	"time"
)

func init() {

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("toml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	viper.AutomaticEnv()

	// 设置日志级别
	log.SetLevel(log.DebugLevel)

	log.SetFormatter(&log.JSONFormatter{})
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
				log.DebugLevel: writer,
				log.InfoLevel:  writer,
				log.WarnLevel:  writer,
				log.ErrorLevel: writer,
				log.FatalLevel: writer,
			}, &log.JSONFormatter{}),
	)

}
func reg() {
	content, err := os.ReadFile("./data/allNews.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	json := string(content)
	value := gjson.Parse(json)

	for _, item := range value.Array() {
		filter.Filter([]byte(item.String()))
	}

	filter.SaveToCSV()
}
func main() {
	tb := pkg.NewTradeBot()
	tb.Run()
	log.Info("running...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmdutil.WaitForSignal(ctx, syscall.SIGINT, syscall.SIGTERM)
}
