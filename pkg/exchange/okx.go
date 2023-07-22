package exchange

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ginvmbot/aitrade/pkg/config"

	"github.com/ginvmbot/newstrade/pkg/fixedpoint"
	"net/http"
	"net/url"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"strings"
	"time"
)

func (e *OkxFuture) InitKline() {
	symbols := e.GetBinanceSymbols()
	strSymbols := make([]string, 0, len(symbols))
	for _, symbol := range symbols {
		fmt.Println(strings.ToLower(symbol.Symbol))
		strSymbols = append(strSymbols, strings.ToLower(symbol.Symbol))
	}
	chunkedSymbols := chunkSlice(strSymbols, 200)
	for _, chunkSymbols := range chunkedSymbols {
		go e.RunWebsocket(chunkSymbols)
		time.Sleep(1 * time.Second)
	}
	//chunkSymbols := []string{
	//	"CFXUSDT",
	//	"IOSTUSDT",
	//	"LITUSDT",
	//}
	//go e.RunWebsocket(chunkSymbols)
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

}

func (e *OkxFuture) GetBinanceSymbols() []config.Symbol {
	resp, err := http.Get("https://fapi.binance.com/fapi/v1/ticker/price")
	if err != nil {
		panic(err) // handle error
	}
	defer resp.Body.Close()

	symbols := make([]config.Symbol, 0, 2300)
	err = json.NewDecoder(resp.Body).Decode(&symbols)
	if err != nil {
		panic(err)
	}

	return symbols
}

func (e *OkxFuture) RunWebsocket(symbols []string) {
	streams := strings.Join(symbols, "@kline_5m/")
	wssUrl, err := url.Parse("wss://fstream.binance.com/stream?streams=" + streams + "@kline_5m")
	fmt.Println(wssUrl)
	if err != nil {
		panic(err)
	}
	opts := websocket.DialOptions{
		CompressionMode: websocket.CompressionDisabled,
	}
	conn, _, err := websocket.Dial(context.Background(), wssUrl.String(), &opts)

	if err != nil {
		panic(err)
	}
	defer conn.Close(websocket.StatusInternalError, "connection closed")

	for {
		//var tickerMsg pkg.TickerMessage
		//if err := wsjson.Read(context.Background(), conn, &tickerMsg); err != nil {
		//	log.Println(err)
		//}
		//
		//ticker := tickerMsg.Data
		//if ticker.Sym == "BTCUSDT" {
		//	fmt.Println(ticker.Sym, ticker)
		//}
		//
		//err = redisDb.SaveTicker(&ticker, ticker.Sym)
		//if err != nil {
		//	log.Println(err)
		//}

		var klineMsg config.KlineMessage
		if err := wsjson.Read(context.Background(), conn, &klineMsg); err != nil {
			fmt.Println("====\n", err)
		}

		kline := klineMsg.Data

		zhengfu := fixedpoint.MustNewFromString(kline.Kline.High).Sub(fixedpoint.MustNewFromString(kline.Kline.Low)).Div(fixedpoint.MustNewFromString(kline.Kline.Low)).Mul(fixedpoint.NewFromFloat(100))

		if zhengfu.Float64() > 0.3 {
			fmt.Println(kline.Symbol, zhengfu)
		}

		//if kline.Symbol == "ETHUSDT" || kline.Symbol == "BTCUSDT" {
		//	fmt.Printf("%s,Open:%s,High:%s,Low:%s,Close:%s,Volume:%s \n", kline.Symbol, kline.Kline.Open, kline.Kline.High, kline.Kline.Low, kline.Kline.Close, kline.Kline.Volume)
		//	fmt.Printf("%s  \n", fixedpoint.MustNewFromString(kline.Kline.Open).Sub(fixedpoint.MustNewFromString(kline.Kline.Close)).Div(fixedpoint.MustNewFromString(kline.Kline.Open)).Percentage())
		//	//fmt.Println(kline.Symbol, kline.Kline.Open, kline.Kline.Close, kline.Kline.High, kline.Kline.Low, kline.Kline.Volume)
		//}

		//err = redisDb.SaveKline(&kline, kline.Symbol)
		//if err != nil {
		//	log.Println(err)
		//}

	}
}
