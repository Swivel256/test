package pkg

import (
	"context"
	"fmt"
	"github.com/ginvmbot/aitrade/pkg/config"
	"github.com/ginvmbot/aitrade/pkg/news"
	server "github.com/ginvmbot/aitrade/pkg/socket"
	"github.com/muesli/cache2go"
	"github.com/tidwall/gjson"
	"regexp"

	"time"
)

type TradeBot struct {
	NewsStream *news.Stream
	RuelEngine *RuelEngine
	//BinanceFuture *exchange.BinanceFuture
	CanOrder chan *config.CanOrder
	Cache    *cache2go.CacheTable
	Twid     *cache2go.CacheTable
	Ctx      context.Context
	Wsserver *server.WebSocketServer
}

func NewTradeBot() TradeBot {
	ctx, _ := context.WithCancel(context.Background())
	cache := cache2go.Cache("Order")
	//exchange := exchange.NewBinanceFuture()
	canorder := make(chan *config.CanOrder)
	ruleEngine := NewRuelEngine(canorder)

	return TradeBot{

		Cache:      cache,
		NewsStream: news.Treenews(ctx),
		RuelEngine: &ruleEngine,
		//BinanceFuture: &exchange,
		CanOrder: canorder,
	}
}

// 启动电报
func (tb *TradeBot) TeleBot() {
	InitBot()
}
func (tb *TradeBot) ListenNews() {
	parseWebSocketEvent := func(message []byte) (interface{}, error) {
		tb.RuelEngine.ParseNews(message)
		return nil, nil
	}

	tb.NewsStream.SetParser(parseWebSocketEvent)
}
func (tb *TradeBot) Run() {
	tb.TeleBot()
	go tb.ListenNews()
	go tb.ListenCreateOrder()
	go tb.SocketServer()

}
func (tb *TradeBot) SocketServer() {
	parseWebSocketEvent := func(message []byte) interface{} {
		value := gjson.ParseBytes(message)
		value.String()

		twid := value.Get("info.twitterId")
		fmt.Println("======twid=========")
		fmt.Println(twid, tb.Twid.Exists(twid.String()))
		if !tb.Twid.Exists(twid.String()) {
			fmt.Println("fuck ==========")
			tb.RuelEngine.ParseNews(message)
		}
		return nil

	}
	ws := server.NewWebSocketServer("0.0.0.0:5555", parseWebSocketEvent)
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
