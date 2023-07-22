// Command remote is a chromedp example demonstrating how to connect to an
// existing Chrome DevTools instance using a remote WebSocket URL.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	devtoolsWsURL := "wss://localhost:9222"
	//var aa, queryNestedSelector string
	flag.Parse()
	if devtoolsWsURL == "" {
		log.Fatal("must specify -devtools-ws-url")
	}

	// create allocator context for use with creating a browser context later
	allocatorContext, cancel := chromedp.NewRemoteAllocator(context.Background(), devtoolsWsURL)
	defer cancel()

	// create context
	ctx, cancel := chromedp.NewContext(allocatorContext)
	defer cancel()
	//
	//ctx, cancel := chromedp.NewContext(
	//	context.Background(),
	//	// chromedp.WithDebugf(log.Printf),
	//)
	//defer cancel()
	ticker := time.NewTicker(10 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println(t)
				ok(ctx)
			}
		}
	}()

	time.Sleep(60 * time.Hour) // 运行1分钟
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")
}
func ok(ctx context.Context) {

	const numWords = 100

	// 生成300个词的文本
	text := Paragraph()

	// 拆分为词
	words := strings.Split(text, " ")

	var output strings.Builder

	// 取前300个词
	for i := 0; i < numWords && i < len(words); i++ {
		output.WriteString(words[i] + " ")
	}
	aa := strings.SplitN(output.String(), " ", 50)
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://twitter.com/home`),
		chromedp.WaitVisible(`[data-testid="SideNav_NewTweet_Button"]`),
		chromedp.Click(`[data-testid="SideNav_NewTweet_Button"]`),
		chromedp.Click(`div[data-testid="tweetTextarea_0"] span[data-offset-key]`),
		chromedp.SendKeys(`div[data-testid="tweetTextarea_0"]  `, strings.Join(aa, "")),

		chromedp.Click(`[data-testid="tweetButton"]`),
	)
	//fmt.Println(aa, queryNestedSelector)

	if err != nil {
		log.Fatal(err)
	}

}
