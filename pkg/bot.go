package pkg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"log"
)

var (
	TgBotApiToken string
	ChatId        int64
	Bot           *tgbotapi.BotAPI
)

type Config struct {
	TgBotApiToken string
	TgChatid      int64
	WeiboUid      []int
	MergeMessage  bool
	Interval      int
	SavePicLocal  bool
	SendLivePics  bool
}

func InitBot() {

	//viper.AddConfigPath(".")
	//if _, file := os.Stat("config.toml"); os.IsNotExist(file) {
	//	viper.SetConfigName("config")
	//	viper.SetConfigType("toml")
	//
	//	viper.SetDefault("TgBotApiToken", "")
	//	viper.SetDefault("TgChatid", "")
	//
	//	if err := viper.SafeWriteConfig(); err != nil {
	//		log.Fatal("保存配置文件失败", err)
	//	}
	//	log.Fatal("根据要求填写 Config.toml 后运行")
	//}
	//
	//if err := viper.ReadInConfig(); err != nil {
	//	log.Fatal("加载配置文件错误", err)
	//}
	//
	//if err := viper.Unmarshal(&config); err != nil {
	//	log.Fatal("解析配置文件错误", err)
	//}
	//fmt.Println(viper.GetString("tgbotapitoken"))

	TgBotApiToken, ChatId = viper.GetString("tgbotapitoken"), viper.GetInt64("tgchatid")

	//fmt.Println(TgBotApiToken, ChatId)
	bot, err := tgbotapi.NewBotAPI(TgBotApiToken)
	if err != nil {
		log.Fatal("连接 Telegram 失败", err)
	}
	Bot = bot

}
