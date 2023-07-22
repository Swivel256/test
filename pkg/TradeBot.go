package pkg

import (
	"bufio"
	"context"
	"fmt"
	"github.com/ginvmbot/aitrade/pkg/config"
	"github.com/ginvmbot/aitrade/pkg/filter"
	news2 "github.com/ginvmbot/aitrade/pkg/maker"
	"github.com/ginvmbot/aitrade/pkg/news"
	server "github.com/ginvmbot/aitrade/pkg/socket"
	twitter2 "github.com/ginvmbot/aitrade/pkg/twitter"
	"github.com/muesli/cache2go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"os"
	"path/filepath"
	"regexp"

	"time"
)

type TradeBot struct {
	NewsStream *news.Stream
	RuelEngine *RuelEngine
	//BinanceFuture *exchange.BinanceFuture
	CanOrder      chan *config.CanOrder
	Cache         *cache2go.CacheTable
	Twid          *cache2go.CacheTable
	Ctx           context.Context
	Wsserver      *server.WebSocketServer
	TwitterClient *twitter2.TwitterClient
}

func NewTradeBot() TradeBot {
	ctx, _ := context.WithCancel(context.Background())
	cache := cache2go.Cache("Order")
	//exchange := exchange.NewBinanceFuture()
	canorder := make(chan *config.CanOrder)
	ruleEngine := NewRuelEngine(canorder)

	return TradeBot{
		TwitterClient: twitter2.NewTwitterClient(),
		Cache:         cache,
		NewsStream:    news.Treenews(ctx),
		RuelEngine:    &ruleEngine,
		//BinanceFuture: &exchange,
		CanOrder: canorder,
	}
}

// 启动电报
func (tb *TradeBot) TwitterFilter() error {

	cache := cache2go.Cache("Order")
	f, err := os.Open(filepath.FromSlash("twid.csv"))
	if err != nil {
		return err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			lines = append(lines, line)
			cache.Add(line, 0, 0)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	fmt.Println("cache.Count()")
	fmt.Println(cache.Count())
	tb.Twid = cache
	return nil
}

// 启动电报
func (tb *TradeBot) TeleBot() {
	InitBot()
}
func (tb *TradeBot) ListenNews() {
	//tb.TwitterFilter()
	All := viper.GetBool("all")
	parseWebSocketEvent := func(message []byte) (interface{}, error) {

		go tb.SaveJson(message)
		maySend := filter.Filter(message)
		if All {
			TeleNews(message)
		} else {
			if maySend {
				log.Info(message,maySend)
				TeleNews(message)
			}
		}

		return nil, nil
	}
	tb.NewsStream.SetParser(parseWebSocketEvent)
}
func (tb *TradeBot) SaveJson(json []byte) {
	value := gjson.ParseBytes(json)

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

	//txt := fmt.Sprintf("%s \n source:%s \n @easytradenews,@pmkooker, @金币瞬间 ", value.Get("title").String(), url)
	//go tb.TwitterClient.PostTwitter(txt)

	log.Info(fmt.Sprintf("title: %s, body: %s, image: %s, url: %s, source: %s, createTime: %s", title, body, image, url, source, createTime))

}

func (tb *TradeBot) ListenNewsMaker() {
	//tb.TwitterFilter()
	//All := viper.GetBool("all")
	parseWebSocketEvent := func(message []byte) (interface{}, error) {

		fmt.Println(string(message))

		return nil, nil
	}
	ctx, _ := context.WithCancel(context.Background())

	makerNews := news2.Makernews(ctx)
	makerNews.SetParser(parseWebSocketEvent)
}
func (tb *TradeBot) Run() {
	tb.TeleBot()
	//go tb.ListenNewsMaker()
	go tb.ListenNews()
	//go tb.ListenCreateOrder()
	go tb.SocketServer()

}
func (tb *TradeBot) SocketServer() {
	parseWebSocketEvent1 := func(message []byte) interface{} {
		//go tb.SaveJson(message)
		maySend := filter.Filter(message)
		fmt.Println("websocket3333", maySend)
		if maySend {
			TeleNews(message)
		}
		return nil

	}
	address := viper.GetString("socket")
	ws := server.NewWebSocketServer(address, parseWebSocketEvent1)
	tb.Wsserver = ws
}
func (tb *TradeBot) ListenCreateOrder() {
	for {
		select {

		case t := <-tb.CanOrder:
			//fmt.Println("fuck", t.PlaceOrder, t.RuelName, t.Info)
			tb.ProcessOrder(t)
		}
	}
	//data := OrderInfo{
	//	Symbol:   i,
	//	Exchange: "binance",
	//}
	//bspot.InitData(data)

}

func (tb *TradeBot) ProcessOrder(t *config.CanOrder) {

	if t.Info["binance-futures"] != "" {
		symbol := t.Info["binance-futures"]
		matched, _ := regexp.MatchString(`(?i).*usdt$`, symbol)
		cached := tb.Cache.Exists(symbol)
		if matched {
			//缓存里不存在，则下单
			if !cached {
				//一个小时内不重复下单
				tb.Cache.Add(symbol, 1*time.Hour, symbol)

				fmt.Println("开合约了：", t.Info["binance-futures"])
			} else {
				fmt.Println("已经有单了，不开了：", t.Info["binance-futures"])

			}
		}
		fmt.Println("缓存数量：", tb.Cache.Count())

	} else if t.Info["binance"] != "" {
		symbol := t.Info["binance"]
		matched, _ := regexp.MatchString(`(?i).*usdt$`, symbol)
		cached := tb.Cache.Exists(symbol)
		if matched {
			//缓存里不存在，则下单
			if !cached {
				//一个小时内不重复下单
				tb.Cache.Add(symbol, 1*time.Hour, symbol)

				fmt.Println("开现货：", symbol)
			} else {
				fmt.Println("已经有单了，不开了：", symbol)

			}
		}
		fmt.Println("缓存数量：", tb.Cache.Count())

	}

	switch t.RuelName {
	case "Upbit":
	case "Coinbase":
	case "Binance_Bnb":
		if t.Info["binance-futures"] != "" {

		}
	}
}
