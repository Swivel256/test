package news

import (
	"context"
	"fmt"
	"github.com/ginvmbot/aitrade/pkg/stream"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"log"
)

//var addr = flag.String("addr", "news.treeofalpha.com", "http service address")

type Stream struct {
	stream.StandardStream
}

func (s *Stream) createEndpoint(ctx context.Context) (string, error) {

	url := "wss://news.treeofalpha.com/ws"
	////url := "ws://127.0.0.1:9999/ws"
	//fmt.Println("url", url)
	return url, nil
}

func (s *Stream) handleDisconnect() {

	s.Reconnect()
}

func (s *Stream) handleConnect() {
	TreeToken := viper.GetString("treetoken")
	t := fmt.Sprintf("login %s", TreeToken)
	fmt.Println(t)
	s.Conn.WriteMessage(websocket.TextMessage, []byte(t))

}

type Event struct {
}

func (s *Stream) dispatchEvent(e interface{}) {

	return

}
func NewStream() *Stream {

	stream := &Stream{
		StandardStream: stream.NewStandardStream(),
	}
	stream.SetDispatcher(stream.dispatchEvent)
	stream.SetDispatcher(stream.dispatchEvent)

	stream.SetEndpointCreator(stream.createEndpoint)

	stream.OnDisconnect(stream.handleDisconnect)
	stream.OnConnect(stream.handleConnect)

	return stream
}

func Treenews(ctx context.Context) *Stream {

	stream := NewStream()

	if err := stream.Connect(ctx); err != nil {
		log.Fatal(err)
	}

	return stream

}

//func (s *Stream) ok() {
//	ticker := time.NewTicker(time.Second)
//	defer ticker.Stop()
//	done := make(chan struct{})
//	interrupt := make(chan os.Signal, 1)
//
//	for {
//		select {
//		case <-done:
//			return
//		case t := <-ticker.C:
//			err := s.Conn.WriteMessage(websocket.TextMessage, []byte(t.String()))
//			if err != nil {
//				log.Println("write:", err)
//				return
//			}
//		case <-interrupt:
//			log.Println("interrupt")
//
//			// Cleanly close the connection by sending a close message and then
//			// waiting (with timeout) for the server to close the connection.
//			err := s.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
//			if err != nil {
//				log.Println("write close:", err)
//				return
//			}
//			select {
//			case <-done:
//			case <-time.After(time.Second):
//			}
//			return
//		}
//	}
//}

//func main() {
//	pkg.InitBot()
//	go treenews()
//	//go newstest()
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	cmdutil.WaitForSignal(ctx, syscall.SIGINT, syscall.SIGTERM)
//}
