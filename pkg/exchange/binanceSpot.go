package exchange

//
//import (
//	"context"
//	"fmt"
//	"github.com/joho/godotenv"
//
//	"github.com/ginvmbot/newstrade/news/newstotelbot/pkg"
//	"net/url"
//	"nhooyr.io/websocket"
//	"nhooyr.io/websocket/wsjson"
//
//	log "github.com/sirupsen/logrus"
//	"github.com/spf13/viper"
//	"github.com/ginvmbot/newstrade/pkg/exchange/binance"
//	"github.com/ginvmbot/newstrade/pkg/fixedpoint"
//	"github.com/ginvmbot/newstrade/pkg/types"
//	"os"
//	"strings"
//)
//
//type BinanceSpot struct {
//	Exchange   *binance.Exchange
//	MinProfit  fixedpoint.Value
//	MaxProfit  fixedpoint.Value
//	StopLoss   fixedpoint.Value
//	StepProfit fixedpoint.Value
//	Market     types.Market
//	Ctx        context.Context
//	Usdt       fixedpoint.Value
//	Symbol     string
//}
//
//func init() {
//	if _, err := os.Stat(".env.local"); err == nil {
//		if err := godotenv.Load(".env.local"); err != nil {
//			log.Fatal(err)
//		}
//	}
//	viper.AutomaticEnv()
//	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
//
//}
//func NewBinanceSpot() BinanceSpot {
//
//	ctx, _ := context.WithCancel(context.Background())
//	//defer cancel()
//
//	key, secret := viper.GetString("binance-api-keys"), viper.GetString("binance-api-secrets")
//
//	fmt.Println(key, secret)
//
//	minProfit, maxProfit, stopLoss, stepProfit := viper.GetFloat64("min-profit"), viper.GetFloat64("max-profit"), viper.GetFloat64("stop-loss"), viper.GetFloat64("step-profit")
//
//	fmt.Println(minProfit, "minProfit", fixedpoint.NewFromFloat(minProfit))
//	var exchange = binance.New(key, secret)
//	//exchange.IsFutures = true
//	return BinanceSpot{
//		Ctx:        ctx,
//		Exchange:   exchange,
//		MinProfit:  fixedpoint.NewFromFloat(minProfit),
//		MaxProfit:  fixedpoint.NewFromFloat(maxProfit),
//		StopLoss:   fixedpoint.NewFromFloat(stopLoss),
//		StepProfit: fixedpoint.NewFromFloat(stepProfit),
//	}
//}
//
//// 预处理数据
//func (e BinanceSpot) InitData(data pkg.OrderInfo) {
//	//go e.TickPrice(data.Symbol)
//
//	e.SubmitOrder()
//	e.Symbol = data.Symbol
//	go e.CheckBalance()
//}
//
//// 检测资产
//func (e BinanceSpot) CheckBalance() {
//	balance, err := e.Exchange.QueryAccountBalances(e.Ctx)
//	e.Usdt = balance["USDT"].Available
//	fmt.Println(balance, err, e.Usdt)
//}
//
//// 获取实时价格
//func (e BinanceSpot) TickPrice(symbol string) {
//
//	wssUrl, err := url.Parse("wss://stream.binance.com/stream?streams=" + strings.ToLower(symbol) + "@kline_1s")
//
//	fmt.Println(wssUrl)
//	if err != nil {
//		panic(err)
//	}
//
//	conn, _, err := websocket.Dial(context.Background(), wssUrl.String(), nil)
//	if err != nil {
//		fmt.Println(err)
//		panic(err)
//	}
//	defer conn.Close(websocket.StatusInternalError, "connection closed")
//
//	for {
//		var klineMsg pkg.KlineMessage
//		if err := wsjson.Read(context.Background(), conn, &klineMsg); err != nil {
//			log.Println(err)
//		}
//
//		kline := klineMsg.Data
//		//if kline.Symbol == symbol {
//		fmt.Println(kline.Symbol, kline.Kline.Open, kline.Kline.Close, kline.Kline.High, kline.Kline.Low, kline.Kline.Volume)
//		//}
//
//		//err = redisDb.SaveKline(&kline, kline.Symbol)
//		//if err != nil {
//		//	log.Println(err)
//		//}
//
//	}
//
//}
//
//// 下单
//func (e BinanceSpot) SubmitOrder() {
//	createdOrder, err := e.Exchange.SubmitOrder(e.Ctx, types.SubmitOrder{
//		Symbol: e.Symbol,
//		Side:   types.SideTypeBuy,
//		Type:   types.OrderTypeMarket,
//		Params: map[string]interface{}{
//			"QuoteOrderQty": e.Usdt,
//		},
//	})
//	print(createdOrder, err)
//
//}
