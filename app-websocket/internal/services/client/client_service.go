package client

import (
	"app-websocket/internal/domain"
	"app-websocket/internal/ports/ws"
	"context"
)

type ClientStorage interface {
	AddClient(ctx context.Context, req *domain.Client) error
}

type MessagePusher interface {
	Produce(msg domain.Event, topic string, key int64) error
}

type ClientService struct {
	cl_storage ClientStorage
	pusher     MessagePusher
	hub        *ws.Hub
}

func NewClientService(clstorage ClientStorage, hub *ws.Hub) *ClientService {
	return &ClientService{
		cl_storage: clstorage,
		hub:        hub,
	}
}

func (r *ClientService) Subscribe(ctx context.Context, client *ws.ClientConnection) error {
	cl := connToDBreq(client)
	//check room password
	if err := r.cl_storage.AddClient(ctx, cl); err != nil {
		return err
	}
	r.hub.AddConnection(client)
	return nil
}

func (r *ClientService) PushMessage(_ context.Context, msg *domain.Message) error {
	return r.pusher.Produce(msg, "messages", msg.RoomID)
}
