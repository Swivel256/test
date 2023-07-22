package html

//
//import (
//	"context"
//	"github.com/ginvmbot/aitrade/pkg/cmdutil"
//	"github.com/ginvmbot/aitrade/pkg/config"
//	exchange2 "github.com/ginvmbot/aitrade/pkg/exchange"
//	"syscall"
//)
//
//func main() {
//	bspot := exchange2.NewOkxFuture()
//
//	Symbols := []string{
//		"TOMOUSDT",
//	}
//	for _, i := range Symbols {
//		data := config.OrderInfo{
//			Symbol:   i,
//			Exchange: "binance",
//		}
//		bspot.InitData(data)
//	}
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	cmdutil.WaitForSignal(ctx, syscall.SIGINT, syscall.SIGTERM)
//
//}
