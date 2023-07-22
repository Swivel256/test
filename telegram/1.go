package main

import (
	"context"
	"fmt"
	"github.com/ginvmbot/aitrade/pkg/cmdutil"
	"syscall"
	"time"

	"golang.org/x/net/proxy"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/dcs"
)

func main() {
	// Dial using SOCKS5 proxy.

	sock5, _ := proxy.SOCKS5("tcp", "127.0.0:15235", &proxy.Auth{
		User:     "",
		Password: "",
	}, proxy.Direct)
	dc := sock5.(proxy.ContextDialer)
	appID := 22896582
	appHash := "585e73c277417d1d4098ff7c12972bbe"
	// Get it from bot father.
	//token := "6591732360:AAEZADHiInf3vO4BYdgs4jXHYUmlUzHUs28"
	// Creating connection.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	client := telegram.NewClient(appID, appHash, telegram.Options{
		Resolver: dcs.Plain(dcs.PlainOptions{
			Dial: dc.DialContext,
		}),
	})

	_ = client.Run(ctx, func(ctx context.Context) error {
		fmt.Println("333333333Started")
		return nil
	})

	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()
	cmdutil.WaitForSignal(ctx1, syscall.SIGINT, syscall.SIGTERM)
}
