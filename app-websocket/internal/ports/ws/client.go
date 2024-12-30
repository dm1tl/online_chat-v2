package ws

import (
	"app-websocket/internal/domain"
	"app-websocket/internal/domain/client"
	"context"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type ServiceChatPusher interface {
	PushMessage(ctx context.Context, msg *domain.Message) error
	Unsubscribe(ctx context.Context, client *client.AddClientReq) error
}

type ClientConnection struct {
	Conn    *websocket.Conn
	Message chan *domain.Message
	Client  *domain.Client
	pusher  ServiceChatPusher
}

func NewClientConnection(conn *websocket.Conn, username string, pass string, uid int64, rid int64, pusher ServiceChatPusher) *ClientConnection {
	return &ClientConnection{
		Conn:    conn,
		Message: make(chan *domain.Message),
		Client: &domain.Client{
			Username: username,
			RoomID:   rid,
			ClientID: uid,
			Password: pass,
		},
		pusher: pusher,
	}
}

func (c *ClientConnection) WriteMessage() {
	defer func() {
		err := c.Conn.Close()
		logrus.Error("writeMessage", err)
	}()
	for {
		message, ok := <-c.Message
		if !ok {
			return
		}
		if err := c.Conn.WriteJSON(message); err != nil {
			logrus.Error("can not send message to client", err)
		}
	}
}

func (c *ClientConnection) ReadMessage(ctx context.Context) {
	defer func() {
		//TODO add unregister method
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Error(err)
			}
			logrus.Error(err)
			break
		}

		msg := &domain.Message{
			Content:     string(m),
			RoomID:      c.Client.RoomID,
			Username:    c.Client.Username,
			UserID:      c.Client.ClientID,
			TimeCreated: time.Now(),
		}

		err = c.pusher.PushMessage(ctx, msg)
		if err != nil {
			logrus.Error("failed to push message:", slog.String("error", err.Error()))
		}
	}
}
