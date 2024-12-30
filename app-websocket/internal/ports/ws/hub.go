package ws

import (
	"app-websocket/internal/domain"
	"context"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type MessageConsumer interface {
	Consume(ctx context.Context, handler domain.MessageHandler) error
}

type Hub struct {
	rooms    map[int64]*HubRoom
	consumer MessageConsumer
	mu       sync.Mutex
}

type HubRoom struct {
	room    *domain.Room
	clients map[int64]*ClientConnection
}

func NewHub(consumer MessageConsumer) *Hub {
	return &Hub{
		rooms:    make(map[int64]*HubRoom),
		consumer: consumer,
	}
}

func (h *Hub) AddConnection(client *ClientConnection) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.rooms[client.Client.RoomID] == nil {
		h.rooms[client.Client.RoomID] = &HubRoom{
			room: &domain.Room{
				ID: client.Client.RoomID,
			},
			clients: make(map[int64]*ClientConnection),
		}
	}
	h.rooms[client.Client.RoomID].clients[client.Client.ClientID] = client
}

func (h *Hub) defaultMessageHandler(msg domain.Message) error {
	h.mu.Lock()
	connections := h.rooms[msg.RoomID].clients
	h.mu.Unlock()
	for _, cl := range connections {
		cl.Message <- &msg
	}
	return nil
}

func (h *Hub) Run(ctx context.Context) {
	attempt := 0
	go func() {
		err := h.consumer.Consume(ctx, h.defaultMessageHandler)
		if err != nil {
			logrus.Error("failed to consume message:", err)

			time.Sleep(expBackoff(attempt))
			attempt++
		}
	}()
	<-ctx.Done()
}

func expBackoff(attempt int) time.Duration {
	maxDelay := 30 * time.Second
	backoff := math.Pow(2, float64(attempt))
	delay := time.Duration(backoff) * time.Second
	if delay > maxDelay {
		delay = maxDelay
	}

	jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
	return delay + jitter
}
