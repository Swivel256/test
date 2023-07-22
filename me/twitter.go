package main

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/ginvmbot/aitrade/pkg/cmdutil"
	twitter2 "github.com/ginvmbot/aitrade/pkg/twitter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
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

	//viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)

}

func main() {

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	done := make(chan bool)
	go func() {
		time.Sleep(30 * time.Minute) // 执行 10 秒后停止
		done <- true
	}()
	var i = 0

	for {
		select {
		case <-done:
			fmt.Println("Done!")
			return
		case t := <-ticker.C:
			twitter := twitter2.NewTwitterClient()
			text := fmt.Sprintf("(%d)--%s", i, gofakeit.Question())
			i++

			twitter.PostTwitter(text)
			fmt.Println("Current time: ", t)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmdutil.WaitForSignal(ctx, syscall.SIGINT, syscall.SIGTERM)
}
