package client

import (
	"app-websocket/internal/domain"
	"app-websocket/internal/ports/ws"
)

func connToDBreq(cl *ws.ClientConnection) *domain.Client {
	return &domain.Client{
		RoomID:   cl.Client.RoomID,
		ClientID: cl.Client.ClientID,
		Username: cl.Client.Username,
		Password: cl.Client.Password,
	}
}
