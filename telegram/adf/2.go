package main

import (
	"fmt"
	"log"
	"time"

	tele "gopkg.in/telebot.v3"
)

func main() {
	pref := tele.Settings{
		Token:  "6591732360:AAEZADHiInf3vO4BYdgs4jXHYUmlUzHUs28",
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/hello", func(c tele.Context) error {
		return c.Send("Hello!")
	})

	b.Handle(tele.OnChannelPost, func(c tele.Context) error {
		// Channel posts only.
		msg := c.Message()
		fmt.Println(msg)
		return nil
	})

	b.Handle("/me", func(c tele.Context) error {
		user := &tele.Chat{
			ID: -1001964184000,
		}

		b.Send(user, "text", &tele.SendOptions{
			// ...
		})
		return nil
		//return c.Send("Hello!")
	})

	b.Start()
}
