package chatService

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024 * 1024
)

type client struct {
	userID string
	ws     *websocket.Conn
	send   chan Event
}

func NewClient(user string, con *websocket.Conn) client {
	return client{
		userID: user,
		ws:     con,
		send:   make(chan Event, maxMessageSize),
	}
}

func (c *client) ReadMe(r *room) {
	defer func() {
		r.unregister <- c
		c.ws.Close()
		println("Readme from client ", c.userID, " closed!!!")
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}

		r.ProcessMessageFromUser(c.userID, message)
	}
}

func (c *client) SendMe(r *room) {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		r.unregister <- c
		c.ws.Close()
		println("SendMe from client ", c.userID, " closed!!!")
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, message)
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *client) write(mt int, message Event) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteJSON(message)
}
