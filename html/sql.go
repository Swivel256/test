package main

import (
	"fmt"
	"github.com/ginvmbot/aitrade/pkg/gojsonq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"regexp"
	"strings"
	"time"
)

type TwitterUser struct {
	UserId int    `json:"userid"`
	Cate   string `json:"cate"`
}
type Twitter struct {
	Source int `json:"source"`
}

func GetTwitter1(id string) string {
	v := []TwitterUser{}
	gojsonq.New().File("data/allNews.json").Select("source").Out(&v)

	fmt.Println(v)
	return v[0].Cate
}

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
	title := "GRAYSCALE: Grayscale's Intentions for btc    etf"

	//matched, _ := regexp.MatchString(`(?i).*(eth.*etf)|(btc.*eft).*`, title)
	a1 := time.Now().UnixNano()

	matched, _ := regexp.MatchString(`(?i).*(eth.*etf)|(btc.*etf).*`, title)
	a2 := time.Now().UnixNano()
	fmt.Println(a2 - a1)

	fmt.Println(matched)

	a1 = time.Now().UnixNano()

	p1 := regexp.MustCompile("(?i)(btc|eth|sec|Bitcoin)").MatchString(title) && regexp.MustCompile("(?i)etf").MatchString(title)
	fmt.Println(p1)
	a2 = time.Now().UnixNano()

	fmt.Println(a2 - a1)
	a1 = time.Now().UnixNano()
	containsBtc := strings.Contains(title, "btc")
	containsEtf := strings.Contains(title, "etf")
	containsEth := strings.Contains(title, "eth")
	containsEft := strings.Contains(title, "eft")

	if (containsBtc && containsEtf) ||
		(containsEth && containsEft) {
		a2 = time.Now().UnixNano()

		fmt.Println(a2 - a1)
	}
}
