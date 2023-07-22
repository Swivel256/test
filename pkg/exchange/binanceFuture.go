package exchange

import (
	"context"
	"fmt"
	"github.com/ginvmbot/newstrade/aibot/pkg/config"
	"github.com/joho/godotenv"
	"net/url"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	"github.com/ginvmbot/newstrade/pkg/exchange/binance"
	"github.com/ginvmbot/newstrade/pkg/fixedpoint"
	"github.com/ginvmbot/newstrade/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type BinanceFuture struct {
	Exchange    *binance.Exchange
	MinProfit   fixedpoint.Value
	MaxProfit   fixedpoint.Value
	StopLoss    fixedpoint.Value
	StepProfit  fixedpoint.Value
	Markets     types.MarketMap
	Market      types.Market
	Ctx         context.Context
	Usdt        fixedpoint.Value
	Symbol      string
	CreateOrder *types.Order
	StopOrder   *types.Order

	PostionPrice fixedpoint.Value
	StopPrice    fixedpoint.Value

	ProfitOrder *types.Order
	ProfitPrice fixedpoint.Value
}

func Init() {
	if _, err := os.Stat("/quant/binanceApi/.env.local"); err == nil {
		if err := godotenv.Load("/quant/binanceApi/.env.local"); err != nil {
			log.Fatal(err)
		}
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

}
func NewBinanceFuture() BinanceFuture {

	ctx, _ := context.WithCancel(context.Background())
	//defer cancel()

	key, secret := viper.GetString("binance-api-key"), viper.GetString("binance-api-secret")

	minProfit, maxProfit, stopLoss, stepProfit := viper.GetFloat64("min-profit"), viper.GetFloat64("max-profit"), viper.GetFloat64("stop-loss"), viper.GetFloat64("step-profit")

	fmt.Println(minProfit, "minProfit", fixedpoint.NewFromFloat(minProfit))
	var exchange = binance.New(key, secret)
	exchange.IsFutures = true
	markets, err := exchange.QueryMarkets(ctx)
	if err != nil {
		fmt.Println(err)
	}

	return BinanceFuture{
		Markets:      markets,
		Ctx:          ctx,
		Exchange:     exchange,
		MinProfit:    fixedpoint.NewFromFloat(minProfit),
		MaxProfit:    fixedpoint.NewFromFloat(maxProfit),
		StopLoss:     fixedpoint.NewFromFloat(stopLoss),
		StepProfit:   fixedpoint.NewFromFloat(stepProfit),
		PostionPrice: fixedpoint.Zero,
		StopPrice:    fixedpoint.Zero,
		ProfitPrice:  fixedpoint.Zero,
	}
}

// 预处理数据
func (e *BinanceFuture) InitData(data config.OrderInfo) {
	//go e.TickPrice(data.Symbol)

	e.Symbol = data.Symbol
	e.QueryMarket()
	e.CheckBalance()
	go e.TickPrice(data.Symbol)

}

// 检测资产
func (e *BinanceFuture) QueryMarket() error {

	market, ok := e.Markets[e.Symbol]
	if !ok {

	}
	e.Market = market

	return nil

}

// 检测资产
func (e *BinanceFuture) CheckBalance() {
	balance, err := e.Exchange.QueryAccountBalances(e.Ctx)
	e.Usdt = balance["USDT"].Available
	fmt.Println("余额", balance, err, e.Usdt)
	e.SubmitOrder()
}

// 获取实时价格
func (e *BinanceFuture) TickPrice(symbol string) {

	wssUrl, err := url.Parse("wss://stream.binance.com/stream?streams=" + strings.ToLower(symbol) + "@ticker")

	fmt.Println(wssUrl)
	if err != nil {
		panic(err)
	}

	conn, _, err := websocket.Dial(context.Background(), wssUrl.String(), nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer conn.Close(websocket.StatusInternalError, "connection closed")

	for {
		var tickerMsg config.TickerMessage
		if err := wsjson.Read(context.Background(), conn, &tickerMsg); err != nil {
			log.Println(err)
			fmt.Println(err)
		}
		//fmt.Println(tickerMsg.Data)

		//StopPrice := fixedpoint.MustNewFromString(tickerMsg.Data.CurrentPrice).Mul(fixedpoint.One.Sub(e.StopLoss)).Round(2, 2)
		//
		//limitPrice := fixedpoint.MustNewFromString(tickerMsg.Data.CurrentPrice).Mul(fixedpoint.One.Add(e.MinProfit))
		//
		//fmt.Printf("现价：%s,止赢：%s,止损:%s", tickerMsg.Data.CurrentPrice, limitPrice, StopPrice)
		e.Process(fixedpoint.MustNewFromString(tickerMsg.Data.CurrentPrice))
		//kline := klineMsg.Data
		////if kline.Symbol == symbol {
		//fmt.Println(kline.Symbol, kline.Kline.Open, kline.Kline.Close, kline.Kline.High, kline.Kline.Low, kline.Kline.Volume)
		////}

		//err = redisDb.SaveKline(&kline, kline.Symbol)
		//if err != nil {
		//	log.Println(err)
		//}

	}

}

// 下单
func (e *BinanceFuture) SubmitOrder() {
	ticker, err := e.Exchange.QueryTicker(e.Ctx, e.Symbol)
	if err != nil {
		fmt.Println(err)
	}
	quantity := e.Usdt.Mul(fixedpoint.NewFromFloat(0.2)).Div(ticker.Sell)
	//quantity := e.Market.MinQuantity
	//spew.Dump(e.Market)
	//fmt.Println("quantity:", quantity, e.Market)
	createdOrder, err := e.Exchange.SubmitOrder(e.Ctx, types.SubmitOrder{
		Market:   e.Market,
		Symbol:   e.Symbol,
		Side:     types.SideTypeBuy,
		Type:     types.OrderTypeMarket,
		Quantity: quantity,
	})

	if err != nil {
		fmt.Println(err)
	}

	e.CreateOrder = createdOrder
	e.PostionPrice = createdOrder.AveragePrice

}

func (e *BinanceFuture) Process(price fixedpoint.Value) {
	//如果没有开止损单，开止损
	//fmt.Println(e.CreateOrder, e.PostionPrice, e.PostionPrice.Compare(fixedpoint.Zero))
	//fmt.Println(e.CreateOrder, e.PostionPrice.Compare(fixedpoint.Zero))

	//fmt.Println(e.PostionPrice.Compare(fixedpoint.Zero))
	fmt.Printf("开仓价：%s,是否为0 %s,止损价：%s \n", e.PostionPrice, price, e.StopPrice)

	//return
	//return
	if e.PostionPrice.IsZero() == true {
		return
	}
	if e.StopPrice.IsZero() == true {
		e.StopPrice = e.PostionPrice.Mul(fixedpoint.One.Sub(e.StopLoss))
		e.ProfitPrice = e.PostionPrice.Mul(fixedpoint.One.Add(e.MinProfit))
		e.SubmitStopOrder(e.StopPrice)
		e.SubmitProfitOrder(e.ProfitPrice)
	} else {
		fmt.Println("追踪止损开始")
		fmt.Printf("开仓价：%s,现价：%s,止损价：%s,价差：%s,实差：%s,比较%s \n",
			e.PostionPrice, price, e.StopPrice, e.PostionPrice.Mul(e.MinProfit), price.Sub(e.StopPrice), price.Sub(e.StopPrice).Compare(e.PostionPrice.Mul(e.MinProfit)))
		//fmt.Println(price, e.StopPrice, price.Sub(e.StopPrice).Compare(e.MinProfit))
		//
		if price.Sub(e.StopPrice).Compare(e.PostionPrice.Mul(e.MinProfit)) == 1 {
			//openOrders, err := e.Exchange.QueryOpenOrders(e.Ctx, e.Symbol)
			//if err != nil {
			//	log.WithError(err).Errorf("query open orders error")
			//} else {
			//	// canceling open orders
			//	fmt.Println(len(openOrders))
			//	if err = e.Exchange.CancelOrders(e.Ctx, openOrders...); err != nil {
			//		log.WithError(err).Errorf("query open orders error")
			//	}
			//}
			//
			//e.StopPrice = price.Mul(fixedpoint.One.Sub(e.StopLoss))
			//
			//e.SubmitStopOrder(e.StopPrice)
		}

	}

}

//func (e *BinanceFuture) ProcessTest(price fixedpoint.Value) {
//	//如果没有开止损单，开止损
//
//	e.PostionPrice = fixedpoint.NewFromFloat(0.3948)
//	e.StopPrice = e.PostionPrice.Mul(fixedpoint.One.Sub(e.StopLoss))
//	e.ProfitPrice = e.PostionPrice.Mul(fixedpoint.One.Add(e.MinProfit))
//	fmt.Println(e.StopPrice, e.ProfitPrice)
//	//limitPrice := createdOrder.AveragePrice.Mul(fixedpoint.One.Add(c.MinProfit))
//	//stopPrice := createdOrder.AveragePrice.Mul(fixedpoint.One.Sub(c.StopLoss))
//	e.SubmitStopOrder(e.StopPrice)
//	e.SubmitProfitOrder(e.ProfitPrice)
//	//
//
//}

// 移动止损
func (e *BinanceFuture) SubmitStopOrder(price fixedpoint.Value) {

	createdStopOrder, err := e.Exchange.SubmitOrder(e.Ctx, types.SubmitOrder{
		Market:     e.Market,
		Symbol:     e.Symbol,
		Side:       types.SideTypeSell,
		Type:       types.OrderTypeStopLimit,
		Price:      price,
		StopPrice:  price,
		Quantity:   e.CreateOrder.Quantity,
		ReduceOnly: true,
	})
	fmt.Println(err)
	e.StopOrder = createdStopOrder
	e.StopPrice = createdStopOrder.Price
}

// 移动止损
func (e *BinanceFuture) SubmitProfitOrder(price fixedpoint.Value) {

	Order, err := e.Exchange.SubmitOrder(e.Ctx, types.SubmitOrder{
		Market:     e.Market,
		Symbol:     e.Symbol,
		Side:       types.SideTypeSell,
		Type:       types.OrderTypeTakeProfitLimit,
		Price:      price,
		StopPrice:  price,
		Quantity:   e.CreateOrder.Quantity,
		ReduceOnly: true,
	})
	fmt.Println(err)
	e.ProfitOrder = Order
	e.ProfitPrice = Order.StopPrice
}
