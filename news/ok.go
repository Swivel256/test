package main

import (
	"context"
	"fmt"
	"github.com/ginvmbot/newstrade/pkg/cmd/cmdutil"
	"github.com/ginvmbot/newstrade/pkg/exchange/binance"
	"github.com/ginvmbot/newstrade/pkg/exchange/okex"
	"github.com/ginvmbot/newstrade/pkg/exchange/okex/okexapi"
	"github.com/ginvmbot/newstrade/pkg/fixedpoint"
	"github.com/ginvmbot/newstrade/pkg/types"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
	"syscall"
)

func init() {
	rootCmd.PersistentFlags().String("okex-api-key", "", "okex api key")
	rootCmd.PersistentFlags().String("okex-api-secret", "", "okex api secret")
	rootCmd.PersistentFlags().String("okex-api-passphrase", "", "okex api secret")
	rootCmd.PersistentFlags().String("symbol", "ETHUSDTSWAP", "symbol")
}

var rootCmd = &cobra.Command{
	Use:   "okex-book",
	Short: "okex book",

	// SilenceUsage is an option to silence usage when an error occurs.
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		symbol := viper.GetString("symbol")
		if len(symbol) == 0 {
			return errors.New("empty symbol")
		}

		key, secret, passphrase := viper.GetString("okex-api-key"),
			viper.GetString("okex-api-secret"),
			viper.GetString("okex-api-passphrase")
		if len(key) == 0 || len(secret) == 0 {
			return errors.New("empty key, secret or passphrase")
		}
		fmt.Println(key, secret, passphrase)

		client := okexapi.NewClient()
		client.Auth(key, secret, passphrase)

		exchange := okex.New(key, secret, passphrase)
		exchange.IsFutures = true
		markets, err := exchange.QueryMarkets(ctx)
		if err != nil {
			return err
		}
		//
		////
		//fmt.Println(symbol, "symbol")
		market, ok := markets[symbol]
		if !ok {
			return fmt.Errorf("market %s is not defined sadfasdf", symbol)
		}
		//
		//fmt.Sprintf("%s", markets)
		//spew.Dump(markets)
		long(exchange, market)

		//
		//instruments, err := client.PublicDataService.NewGetInstrumentsRequest().
		//	InstrumentType("SPOT").Do(ctx)
		//if err != nil {
		//	return err
		//}
		//
		//log.Infof("instruments: %+v", instruments)
		//
		//fundingRate, err := client.PublicDataService.NewGetFundingRate().InstrumentID("BTC-USDT-SWAP").Do(ctx)
		//if err != nil {
		//	return err
		//}
		//log.Infof("funding rate: %+v", fundingRate)
		//
		//log.Infof("ACCOUNT BALANCES:")
		//account, err := client.AccountBalances()
		//if err != nil {
		//	return err
		//}
		//
		//log.Infof("%+v", account)
		//
		//log.Infof("ASSET BALANCES:")
		//assetBalances, err := client.AssetBalances()
		//if err != nil {
		//	return err
		//}
		//
		//for _, balance := range assetBalances {
		//	log.Infof("%T%+v", balance, balance)
		//}
		//
		//log.Infof("ASSET CURRENCIES:")
		//currencies, err := client.AssetCurrencies()
		//if err != nil {
		//	return err
		//}
		//
		//for _, currency := range currencies {
		//	log.Infof("%T%+v", currency, currency)
		//}
		//
		//log.Infof("MARKET TICKERS:")
		//tickers, err := client.MarketTickers(okexapi.InstrumentTypeSpot)
		//if err != nil {
		//	return err
		//}
		//
		//for _, ticker := range tickers {
		//	log.Infof("%T%+v", ticker, ticker)
		//}
		//
		//ticker, err := client.MarketTicker("ETH-USDT")
		//if err != nil {
		//	return err
		//}
		//log.Infof("TICKER:")
		//log.Infof("%T%+v", ticker, ticker)

		//createdOrder, err := exchange.SubmitOrder(ctx, types.SubmitOrder{
		//	Symbol: "ETHUSDT",
		//	Market: market,
		//	Side:   types.SideTypeBuy,
		//	Type:   types.OrderTypeMarket,
		//	Price:  fixedpoint.NewFromFloat(1300),
		//	//StopPrice:        fixedpoint.NewFromFloat(price),
		//	Quantity: fixedpoint.NewFromFloat(0.1),
		//
		//	//TimeInForce: "GTC",
		//	//ReduceOnly:       true,
		//})
		//if err != nil {
		//	fmt.Println(err)
		//}
		//
		//fmt.Println("create----:\n", market)

		//go long(exchange, market)
		//postion, _ := client.AccountPositions()
		//fmt.Println(string(postion))
		//data, _ := client.ClosePositions("long")
		////
		//fmt.Println("----:", data)
		//

		//algoId, _ := client.GetPendingAlgos()
		////	orderid := gjson.Get(data1.String(), "data[0].algoId")
		////
		//fmt.Println("fuck 止损单：", algoId)
		//
		//client.CancelAlgos(algoId)
		////
		//go stopLong(exchange, market)
		//
		//createdOrder, err := exchange.SubmitOrder(ctx, types.SubmitOrder{
		//	Symbol: "ETHUSDTSWAP",
		//	Market: market,
		//	Side:   types.SideTypeSell,
		//	Type:   types.OrderTypeMarket,
		//	//Price:  fixedpoint.NewFromFloat(price),
		//	//StopPrice:        fixedpoint.NewFromFloat(price),
		//	Quantity: fixedpoint.NewFromFloat(1),
		//
		//	TimeInForce: "GTC",
		//	ReduceOnly:  true,
		//})
		//if err != nil {
		//	fmt.Println(err)
		//}
		//
		//log.Info(createdOrder)

		cmdutil.WaitForSignal(ctx, syscall.SIGINT, syscall.SIGTERM)
		return nil
	},
}

// // 开多自带止赢止损
//
//	func long(exchange *okex.Exchange, market types.Market) {
//		ctx, cancel := context.WithCancel(context.Background())
//		defer cancel()
//		price := 1700.0
//		quantity := 1.0
//	}
//
// 开多自带止赢止损
func long(exchange *okex.Exchange, market types.Market) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	price := 1920.0
	quantity := 1.0
	//市价单买入
	createdOrder, err := exchange.SubmitOrder(ctx, types.SubmitOrder{
		Symbol:    "ETHUSDTSWAP",
		Market:    market,
		Side:      types.SideTypeBuy,
		Type:      types.OrderTypeMarket,
		Price:     fixedpoint.NewFromFloat(price),
		StopPrice: fixedpoint.NewFromFloat(price),
		Quantity:  fixedpoint.NewFromFloat(quantity),

		TimeInForce: "GTC",
		//ReduceOnly:       true,
	})
	if err != nil {
		fmt.Println(err)
	}
	//
	fmt.Println("create----:\n", createdOrder, err)
	//exchange.CancelOrders(ctx)
	//限价单止赢
	//
	//createdTackProfitOrder, err := exchange.SubmitOrder(ctx, types.SubmitOrder{
	//	Symbol:    "ETHUSDTSWAP",
	//	Market:    market,
	//	Side:      types.SideTypeSell,
	//	Type:      types.OrderTypeTakeProfitLimit,
	//	Price:     fixedpoint.NewFromFloat(1840),
	//	StopPrice: fixedpoint.NewFromFloat(1840),
	//	Quantity:  fixedpoint.NewFromFloat(1),
	//
	//	TimeInForce: "GTC",
	//	ReduceOnly:  true,
	//})
	//fmt.Println("take,profit----::\n", createdTackProfitOrder, err)
	//限价单止损

	//createdStopOrder, err := exchange.SubmitOrder(ctx, types.SubmitOrder{
	//	Symbol:    "ETHUSDTSWAP",
	//	Market:    market,
	//	Side:      types.SideTypeSell,
	//	Type:      types.OrderTypeStopLimit,
	//	Price:     fixedpoint.NewFromFloat(1680),
	//	StopPrice: fixedpoint.NewFromFloat(1680),
	//	Quantity:  fixedpoint.NewFromFloat(quantity * 10),
	//
	//	TimeInForce: "GTC",
	//	ReduceOnly:  true,
	//})
	//fmt.Println("stop-----:", createdStopOrder, err)

}

// 开多自带止赢止损
func stopLong(exchange *okex.Exchange, market types.Market) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quantity := 0.1

	createdStopOrder, err := exchange.SubmitOrder(ctx, types.SubmitOrder{
		Symbol:    "ETHUSDTSWAP",
		Market:    market,
		Side:      types.SideTypeBuy,
		Type:      types.OrderTypeTakeProfitLimit,
		Price:     fixedpoint.NewFromFloat(1640),
		StopPrice: fixedpoint.NewFromFloat(1640),
		Quantity:  fixedpoint.NewFromFloat(quantity),

		TimeInForce: "GTC",
		ReduceOnly:  true,
	})
	fmt.Println("stop-----:", createdStopOrder, err)

}

// 开空自带止赢止损
func short(exchange *binance.Exchange, market types.Market) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	price := 1340.0
	quantity := 0.9
	//市价单买入
	createdOrder, err := exchange.SubmitOrder(ctx, types.SubmitOrder{
		Symbol: "ETHUSDT",
		Market: market,
		Side:   types.SideTypeSell,
		Type:   types.OrderTypeMarket,
		Price:  fixedpoint.NewFromFloat(price),
		//StopPrice:        fixedpoint.NewFromFloat(price),
		Quantity: fixedpoint.NewFromFloat(quantity),

		TimeInForce: "GTC",
		//ReduceOnly:       true,
	})
	if err != nil {
		fmt.Println(err)
	}

	log.Info(createdOrder)

	//限价单止赢

	createdTackProfitOrder, err := exchange.SubmitOrder(ctx, types.SubmitOrder{
		Symbol:    "ETHUSDT",
		Market:    market,
		Side:      types.SideTypeBuy,
		Type:      types.OrderTypeLimit,
		Price:     fixedpoint.NewFromFloat(price - 20),
		StopPrice: fixedpoint.NewFromFloat(price - 20),
		Quantity:  fixedpoint.NewFromFloat(quantity),

		TimeInForce: "GTC",
		ReduceOnly:  true,
	})
	fmt.Println(createdTackProfitOrder)
	//限价单止损
	//
	//createdStopOrder, err := exchange.SubmitOrder(ctx, types.SubmitOrder{
	//	Symbol:    "ETHUSDT",
	//	Market:    market,
	//	Side:      types.SideTypeBuy,
	//	Type:      types.OrderTypeStopLimit,
	//	Price:     fixedpoint.NewFromFloat(price + 20),
	//	StopPrice: fixedpoint.NewFromFloat(price + 20),
	//	Quantity:  fixedpoint.NewFromFloat(quantity),
	//
	//	TimeInForce: "GTC",
	//	ReduceOnly:  true,
	//})
	//fmt.Println(createdStopOrder)
}

func main() {
	if _, err := os.Stat(".env.local"); err == nil {
		if err := godotenv.Load(".env.local"); err != nil {
			log.Fatal(err)
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.WithError(err).Error("bind pflags error")
	}

	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		log.WithError(err).Error("cmd error")
	}
}
