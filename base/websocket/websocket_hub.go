package websocket

import (
	"encoding/json"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string]*Client

	// Inbound messages from the clients.
	send chan *Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister    chan *Client
	MessageHanler func(c *Client, msg *Message) error
}

//设置一个全局变量
var WSHub *Hub

func NewHub() *Hub {
	WSHub = &Hub{
		send:       make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
	}
	return WSHub
}

func InitWebSocketHub() {
	NewHub()
}

//暂时不能执行两次，可加判断，让函数可重入
func (h *Hub) Run() {

	go func() {
		for {
			select {
			case client := <-h.register:
				h.clients[client.clientId] = client
			case client := <-h.unregister:
				if _, ok := h.clients[client.clientId]; ok {
					delete(h.clients, client.clientId)
					close(client.send)
				}
			case message := <-h.send:
				msg, err := json.Marshal(message)
				if err == nil {
					switch message.MsgType {
					case WS_P2P:
						c := h.clients[message.To]
						if c != nil {
							c.send <- msg
						}

					case WS_BROCAST:
						for _, client := range h.clients {
							select {
							case client.send <- msg:
							default:
								close(client.send)
								delete(h.clients, client.clientId)
							}
						}
					}
				}

			}
		}
	}()

}

func (h *Hub) Send(msg *Message) {
	if msg == nil {
		return
	}
	h.send <- msg
}
