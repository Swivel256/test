package filter

import (
	"fmt"
	gt "github.com/bas24/googletranslatefree"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"strings"
	"time"
)

// 是否发送
func Filter(json []byte) bool {
	data := gjson.ParseBytes(json)
	title := data.Get("title").String()

	if data.Get("en").Exists() {
		title = data.Get("en").String()

	}
	source := data.Get("source").String()

	t := data.Get("time").Int()
	unixTimeStamp := time.Unix(t/1000, 0)
	createTime := unixTimeStamp.Format("2006-01-02 15:04:05")

	var maySend bool

	switch source {
	case "Blogs":
		maySend = BlogFilter(title)
		break
	case "Twitter":

		maySend = TwitterFilter(json)
		break

	case "Upbit":
		maySend = Upbit(title)
		break
	case "Coinbase":
		maySend = Coinbase(title)
		break
	case "Binance EN":
		maySend = BinanceEN(title)
		break

	}
	fmt.Printf("来消息了: %s (%s),%t \n", createTime, title, maySend)
	log.Info(fmt.Sprintf("来消息了: %s (%s),%t \n", createTime, title, maySend))

	return maySend
	//if maySend {
	//	TeleNews(json)
	//} else {
	//	fmt.Println("不符合推送标准")
	//}
}

func TeleNews(json []byte) {

	value := gjson.ParseBytes(json)

	if value.Get("title").Exists() {
		title := value.Get("title").String()
		fmt.Println("===title===", title)
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
		twid := data.Get("info.twitterId").String()

		SendMsg(title, createTime, url, body, source, coins, twid)

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
func Translate(text string) string {
	result, err := gt.Translate(text, "en", "zh")
	//fmt.Println(result)
	if err != nil {
		fmt.Println("google", err)
		return ""
	}
	return result
}
func SendMsg(title, createTime, url, body, source, coins, twid string) {

	Titles = append(Titles, title)
	Time = append(Time, createTime)
	Source = append(Source, source)
	Url = append(Url, url)
	Twid = append(Twid, twid)
	return

	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("*%s*\n", title))
	b.WriteString("————————————\n")
	if len(coins) > 0 {
		b.WriteString(fmt.Sprintf("相关币：%s\n", coins))
	}
	if body != "" {
		b.WriteString(fmt.Sprintf("%s*\n", body))
	}

	b.WriteString(fmt.Sprintf("%s [%s](%s)\n", source, url, url))

	b.WriteString(fmt.Sprintf("时间:%s \n", createTime))
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
		fmt.Println("telegram bot send", err)
	}
}
