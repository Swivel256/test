package main

import (
	"context"
	"fmt"
	"github.com/ginvmbot/aitrade/pkg"
	"github.com/ginvmbot/aitrade/pkg/cmdutil"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"syscall"
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
	tb := pkg.NewTradeBot()
	tb.Run()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmdutil.WaitForSignal(ctx, syscall.SIGINT, syscall.SIGTERM)
}
