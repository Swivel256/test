package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有Origin
	},
}

type WebSocketServer struct {
	*websocket.Conn
	Parse   func([]byte) interface{}
	Address string
}
func NewWebSocketServer(address string, parse func([]byte) interface{}) *WebSocketServer {
	ws:=&WebSocketServer{

		Address: address,
		Parse: parse,
	}
	ws.StartWebSocketServer()
	return ws
}
func (e *WebSocketServer)StartWebSocketServer() {
	http.HandleFunc("/ws", e.handleWebSocketConnection)
	fmt.Println("WebSocket server is listening on", e.Address)
	err := http.ListenAndServe(e.Address, nil)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
}

func (e *WebSocketServer)handleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade connection:", err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			// Check if client has closed connection
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("WebSocket connection closed by client:", err)
			} else {
				fmt.Println("Failed to read message from WebSocket:", err)
			}
			return
		}
		e.Parse(message)
		fmt.Println("Received message:", string(message))

		err = conn.WriteMessage(messageType, message)
		if err != nil {
			fmt.Println("Failed to write message:", err)
			return
		}
	}
}