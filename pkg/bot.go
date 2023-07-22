package pkg

import (
	"fmt"
	gtranslate "github.com/gilang-as/google-translate"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"

	"strings"
	"time"
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

	TgBotApiToken, ChatId = viper.GetString("tgbotapitoken"), viper.GetInt64("tgchatid")
	fmt.Println(TgBotApiToken, ChatId)
	bot, err := tgbotapi.NewBotAPI(TgBotApiToken)
	if err != nil {
		log.Fatal("连接 Telegram 失败", err)
	}
	Bot = bot

}
func Translate(text string) string {
	//fmt.Println("deeprrr", text)
	//title := DeepTranslate(text)
	//fmt.Println("deep", title)
	//if title != "" {
	//
	//	fmt.Println("google", title)
	//
	//	title = GoogleTranslate(title)
	//}

	value := gtranslate.Translate{
		Text: text,
		From: "en",
		To:   "zh",
	}
	translated, err := gtranslate.Translator(value)
	if err != nil {
		log.Error(err)
	} else {
		return translated.Text
	}
	return text

	//title := GoogleTranslate(text)
	////fmt.Println(result)
	//return title

}
func TeleNews(json []byte) {

	value := gjson.ParseBytes(json)
	log.Info("开始发电报", value.String())
	if value.Get("title").Exists() {
		title := value.Get("title").String()

		url := "http://127.0.0.1"

		body := ""
		image := ""
		media := []string{}
		source := ""

		if value.Get("en").Exists() {
			title = value.Get("en").String()
		}
		if value.Get("body").Exists() {
			body = value.Get("body").String()
		}
		if value.Get("image").Exists() {
			image = value.Get("image").String()
			if image != "" {
				media = append(media, image)
			}
		}
		if value.Get("url").Exists() {
			url = value.Get("url").String()
		}
		if value.Get("link").Exists() {
			url = value.Get("link").String()
		}
		if value.Get("source").Exists() {
			source = value.Get("source").String()
		}
		t := value.Get("time").Int()
		unixTimeStamp := time.Unix(t/1000, 0)
		createTime := unixTimeStamp.Format("2006-01-02 15:04:05")
		en := title
		if title != "" && len(title) > 0 {
			title = Translate(title)
		}
		if body != "" && len(body) > 0 {
			body = Translate(body)
		}

		data := value.Get("suggestions.#.coin")
		coins := ""
		if data.Exists() {
			coins = data.String()
		}

		go SendMsg(en, title, createTime, url, body, source, coins)

		//暂时先去掉
		//遍历 e.Keys
		//for _, key := range e.Keys {
		//	canPush := regexp.MustCompile(fmt.Sprintf("(?i)%s", key)).MatchString(title)
		//
		//	if canPush {
		//		SendMessage(title, url, body, media)
		//	}
		//}

	}
}

func SendMsg(en, title, createTime, url, body, source, coins string) {

	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("*%s*\n", en))
	b.WriteString(fmt.Sprintf("*%s*\n", title))
	b.WriteString("\n")
	if body != "" {
		b.WriteString(fmt.Sprintf("*%s*\n", body))
	}

	b.WriteString("\n\n")

	if len(coins) > 0 {
		b.WriteString(fmt.Sprintf("相关币：%s\n", coins))
	}
	b.WriteString("——————————————\n")
	b.WriteString(fmt.Sprintf("时间:%s \n", createTime))

	b.WriteString(fmt.Sprintf("%s:\n [%s](%s)\n", source, url, url))
	fmt.Println("==", ChatId)
	msg := tgbotapi.NewMessage(ChatId, b.String())
	messageInlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("🔗点击查看原文", url),
		),
	)
	msg.ReplyMarkup = messageInlineKeyboard
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.DisableWebPagePreview = false
	_, err := Bot.Send(msg)
	if err != nil {
		log.Error(err)
		fmt.Println("telegram  33333 bot send", err)
	}
}
