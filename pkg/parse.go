package pkg

import (
	"fmt"
	hashmap "github.com/duke-git/lancet/v2/datastructure/hashmap"
	"github.com/ginvmbot/aitrade/pkg/config"
	"github.com/ginvmbot/aitrade/pkg/news"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"os"
	"strings"
	"time"
)

type RuelEngine struct {
	Keys []string
	RuleEngine
	CanOrder      chan *config.CanOrder
	Hashmap       *hashmap.HashMap

}

func NewRuelEngine(co chan *config.CanOrder) RuelEngine {
	content, err := os.ReadFile("./key.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	json := string(content)
	value := gjson.Parse(json)
	engine := RuelEngine{
		CanOrder: co,
	}
	for _, item := range value.Get("keys").Array() {
		engine.Keys = append(engine.Keys, item.String())
	}

	return engine

}
func (e RuelEngine) ParseNews(json []byte) {
	e.OrdreNews(json)

}
func (e RuelEngine) OrdreNews(json []byte) {
	canTrade, name := e.processNews(gjson.ParseBytes(json))
	value := gjson.ParseBytes(json)
	data := value.Get("suggestions.0.symbols")
	//hm := hashmap.NewHashMapWithCapacity(uint64(100), uint64(1000))
	hm := make(map[string]string)

	for _, v := range data.Array() {
		//fmt.Println(v.Get("symbol").String(), v.Get("exchange").String())
		hm[v.Get("exchange").String()] = v.Get("symbol").String()
	}
	e.CanOrder <- &config.CanOrder{

		PlaceOrder: canTrade,
		RuelName:   name,
		Info:       hm,
	}

}

func (e RuelEngine) TeleNews(json []byte) {

	value := gjson.ParseBytes(json)
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

		if title != "" {
			title = news.Translate(title)
		}
		if body != "" {
			body = news.Translate(body)
		}

		data := value.Get("suggestions.#.coin")
		coins := ""
		if data.Exists() {
			coins = data.String()
		}

		SendMsg(title, createTime, url, body, source, coins)


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

func SendMsg(title, createTime, url, body, source, coins string) {
	fmt.Println("====fuck==")

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
	//_, err := Bot.Send(msg)
	//if err != nil {
	//	fmt.Println(err)
	//}
}
