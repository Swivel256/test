// Binary bot-auth-manual implements example of custom session storage and
// manually setting up client options without environment variables.
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gotd/td/telegram/dcs"
	"go.uber.org/zap"
	"golang.org/x/net/proxy"
	"sync"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
)

// Run runs f callback with context and logger, panics on error.
func Run(f func(ctx context.Context, log *zap.Logger) error) {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer func() { _ = log.Sync() }()
	// No graceful shutdown.
	ctx := context.Background()
	if err := f(ctx, log); err != nil {
		log.Fatal("Run failed", zap.Error(err))
	}
	// Done.
}

// memorySession implements in-memory session storage.
// Goroutine-safe.
type memorySession struct {
	mux  sync.RWMutex
	data []byte
}

// LoadSession loads session from memory.
func (s *memorySession) LoadSession(context.Context) ([]byte, error) {
	if s == nil {
		return nil, session.ErrNotFound
	}

	s.mux.RLock()
	defer s.mux.RUnlock()

	if len(s.data) == 0 {
		return nil, session.ErrNotFound
	}

	cpy := append([]byte(nil), s.data...)

	return cpy, nil
}

// StoreSession stores session to memory.
func (s *memorySession) StoreSession(ctx context.Context, data []byte) error {
	s.mux.Lock()
	s.data = data
	s.mux.Unlock()
	return nil
}

func main() {
	// Grab those from https://my.telegram.org/apps.

	appID := 22896582
	appHash := "585e73c277417d1d4098ff7c12972bbe"
	// Get it from bot father.
	token := "6591732360:AAEZADHiInf3vO4BYdgs4jXHYUmlUzHUs28"
	flag.Parse()

	// Using custom session storage.
	// You can save session to database, e.g. Redis, MongoDB or postgres.
	// See memorySession for implementation details.
	sessionStorage := &memorySession{}

	Run(func(ctx context.Context, log *zap.Logger) error {
		sock5, _ := proxy.SOCKS5("tcp", "127.0.0:15235", &proxy.Auth{
			User:     "",
			Password: "",
		}, proxy.Direct)
		dc := sock5.(proxy.ContextDialer)
		client := telegram.NewClient(appID, appHash, telegram.Options{
			SessionStorage: sessionStorage,
			Logger:         log,
			Resolver: dcs.Plain(dcs.PlainOptions{
				Dial: dc.DialContext,
			}),
		})

		return client.Run(ctx, func(ctx context.Context) error {
			// Checking auth status.
			status, err := client.Auth().Status(ctx)
			fmt.Println(err)
			if err != nil {
				return err
			}
			// Can be already authenticated if we have valid session in
			// session storage.
			if !status.Authorized {
				// Otherwise, perform bot authentication.
				if _, err := client.Auth().Bot(ctx, token); err != nil {
					fmt.Println(err)
					return err
				}
			}
			fmt.Println("done")
			// All good, manually authenticated.
			log.Info("Done")

			return nil
		})
	})
}
