package clients

import (
	"OnlineGame/database"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	User *database.User
	Conn *websocket.Conn
	send chan []byte
}

func NewClient(user *database.User, conn *websocket.Conn) *Client {
	return &Client{
		User: user,
		Conn: conn,
		send: make(chan []byte, 256),
	}
}

func (c *Client) ReadLoop(
	onMessage func(userID uint, data []byte),
) {
	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("client %d disconnected: %v", c.User.ID, err)
			return
		}
		onMessage(c.User.ID, data)
	}
}

func (c *Client) WriteLoop() {
	for data := range c.send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return
		}
	}
}

func (c *Client) Send(v any) {
	data, err := json.Marshal(v)
	if err != nil {
		return
	}
	c.send <- data
}

func (c *Client) UpdateConnection(newConn *websocket.Conn) {
	c.Conn = newConn
}
