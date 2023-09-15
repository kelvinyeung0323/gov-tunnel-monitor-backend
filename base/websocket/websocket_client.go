package websocket

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WSMsgType string

var (
	WS_BROCAST WSMsgType = "broadcast" //广播类型
	WS_P2P     WSMsgType = "p2p"       //个人对个人
	WS_CMD     WSMsgType = "cmd"       //命令请求
)

//WebSocket 消息体
type Message struct {
	MsgType WSMsgType `json:"msgType"`
	From    string    `json:"from"`
	To      string    `json:"to"`
	Body    any       `json:"body"`
	Cmd     string    `json:"cmd"`
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	//允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub      *Hub
	clientId string //用户Id
	ctx      *gin.Context
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		//将换行替换成空格
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		//转换为Message
		msg := &Message{}
		err = json.Unmarshal(message, &msg)
		if err == nil {
			msg.From = c.clientId
			c.handlMessage(msg)

		} else {
			log.Printf("error: unmarshal websocket message error:%v", err)
		}

	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C: //超时
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handlMessage(msg *Message) {
	log.Printf("websocket recieved msg: %v", msg)
	// msg.MsgType = WS_BROCAST
	// c.hub.send <- msg
	if c.hub.MessageHanler != nil {
		c.hub.MessageHanler(c, msg)
	}
}

// serveWs handles websocket requests from the peer.
func WebsocketHandler(hub *Hub, ctx *gin.Context) {
	//TODO:获取请求ID,设置clientId

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{ctx: ctx, hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	//把事务管理注入

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()

}
