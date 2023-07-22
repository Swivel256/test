package main

import (
	"context"
	"fmt"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/html"
	"github.com/gotd/td/telegram/uploader"
	"github.com/gotd/td/tg"
	log "github.com/sirupsen/logrus"

	"github.com/gotd/td/telegram"
)

const api_id = 26875732
const api_hash = "3e9b6bccde2480cf6823ae785710f462"

func main() {
	// https://core.telegram.org/api/obtaining_api_id
	client := telegram.NewClient(api_id, api_hash, telegram.Options{})
	fmt.Println(4444)
	if err := client.Run(context.Background(), func(ctx context.Context) error {
		fmt.Println(333)
		api := tg.NewClient(client)

		// Helper for uploading. Automatically uses big file upload when needed.
		u := uploader.NewUploader(api)

		// Helper for sending messages.
		sender := message.NewSender(api).WithUploader(u)
		filePath := "640x360.png"
		// Uploading directly from path. Note that you can do it from
		// io.Reader or buffer, see From* methods of uploader.
		log.Info("Uploading file")
		upload, err := u.FromPath(ctx, filePath)
		if err != nil {
			return fmt.Errorf("upload %q: %w", filePath, err)
		}

		// Now we have uploaded file handle, sending it as styled message.
		// First, preparing message.
		document := message.UploadedDocument(upload,
			html.String(nil, `Upload: <b>From bot</b>`),
		)

		// You can set MIME type, send file as video or audio by using
		// document builder:
		document.
			MIME("audio/mp3").
			Filename("some-audio.mp3").
			Audio()

		// Resolving target. Can be telephone number or @nickname of user,
		// group or channel.
		targetDomain := "."
		target := sender.Resolve(targetDomain)

		// Sending message with media.
		log.Info("Sending file")
		if _, err := target.Media(ctx, document); err != nil {
			return fmt.Errorf("send: %w", err)
		}

		// Return to close client connection and free up resources.
		return nil
	}); err != nil {
		fmt.Println(err)
	}
	fmt.Println(5555)
	// Client is closed.
}
