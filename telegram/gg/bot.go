package main

import (
	"context"
	"fmt"
	"github.com/ginvmbot/aitrade/pkg/cmdutil"
	"log"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	//tgbotapi.NewBotAPIWithClient()
	bot, err := tgbotapi.NewBotAPI("6591732360:AAEZADHiInf3vO4BYdgs4jXHYUmlUzHUs28")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		fmt.Println(update)
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmdutil.WaitForSignal(ctx, syscall.SIGINT, syscall.SIGTERM)
}
