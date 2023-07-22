package news

import (
	"context"
	"fmt"
	"github.com/ginvmbot/aitrade/pkg/stream"
	"log"
)

//var addr = flag.String("addr", "news.treeofalpha.com", "http service address")

type Stream struct {
	stream.StandardStream
}

func (s *Stream) createEndpoint(ctx context.Context) (string, error) {
	url := "wss://db.tuleep.trade/realtime/v1/websocket?apikey=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Im90aW5waXF1Z2Jqa3V1ZWh0d2NxIiwicm9sZSI6ImFub24iLCJpYXQiOjE2NzYwNDcyODcsImV4cCI6MTk5MTYyMzI4N30.O9ZMyq5s5UukgKYYSmE2eqTLwfu8rrwy2t_sij8ZQZQ&vsn=1.0.0"
	return url, nil
}

func (s *Stream) handleDisconnect() {
	fmt.Println("已断开")
	s.Connected = false
	s.Reconnect("已断开")
}

func (s *Stream) handleConnect() {
	fmt.Println("已连接")
	//s.Conn.WriteMessage(websocket.TextMessage, []byte(t))

}

type Event struct {
}

func (s *Stream) dispatchEvent(e interface{}) {

	return

}
func MakerStream() *Stream {

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

func Makernews(ctx context.Context) *Stream {

	stream := MakerStream()

	if err := stream.Connect(ctx); err != nil {
		log.Fatal(err)
	}

	return stream

}
