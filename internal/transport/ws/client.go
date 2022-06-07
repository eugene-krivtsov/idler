package ws

import (
	"context"
	"github.com/eugene-krivtsov/idler/internal/model/entity"
	"github.com/eugene-krivtsov/idler/internal/service"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"time"
)

type Client struct {
	user           string
	pool           *Pool
	conn           *websocket.Conn
	send           chan entity.Message
	messageService service.Messages
}

func NewClient(user string, conn *websocket.Conn, pool *Pool, messageService service.Messages) *Client {
	client := &Client{
		user:           user,
		pool:           pool,
		conn:           conn,
		send:           make(chan entity.Message),
		messageService: messageService,
	}
	client.pool.register <- client

	return client
}

func (c *Client) HandleMessage() {
	defer func() {
		c.pool.unregister <- c
		if err := c.conn.Close(); err != nil {
			logrus.Errorf("error occured on web socket client close: %s", err.Error())
			return
		}
	}()

	for {
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Errorf("unexpected error occured on web socket client close: %s", err.Error())
				return
			}
		}

		message := entity.Message{
			Sender: "a",
			SentAt: time.Now(),
			Text:   string(p),
		}

		c.pool.broadcast <- message

		if err := c.messageService.Save(context.Background(), message); err != nil {
			logrus.Errorf("error occured on web socket sending message: %s", err.Error())
			return
		}
	}
}
