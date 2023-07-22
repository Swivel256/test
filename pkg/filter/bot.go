package filter

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

func init() {

	////TgBotApiToken, ChatId = viper.GetString("tgbotapitoken"), viper.GetInt64("tgchatid")
	////fmt.Println(TgBotApiToken, ChatId)
	//
	//TgBotApiToken = "6378721935:AAGV8ELFAKxnNqNuwYjqs910KgO2Q9LPz1E"
	//ChatId = -982559994
	//bot, err := tgbotapi.NewBotAPI(TgBotApiToken)
	//if err != nil {
	//	log.Fatal("连接 Telegram 失败", err)
	//}
	//Bot = bot

}
