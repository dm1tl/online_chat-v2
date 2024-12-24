package room

import (
	"app-websocket/internal/domain"
	"app-websocket/internal/domain/room"
	"context"
)

const (
	roomTopic = "room-topic"
	roomKey   = int64(1)
)

type RoomPusher interface {
	Produce(msg domain.Event, topic string, key int64) error
}

type RoomStorage interface {
	CreateRoom(ctx context.Context, req room.CreateRoomReq) (int64, error)
}

type RoomService struct {
	pusher  RoomPusher
	storage RoomStorage
}

func NewRoomService(pusher RoomPusher, storage RoomStorage) *RoomService {
	return &RoomService{
		pusher:  pusher,
		storage: storage,
	}
}

func (r *RoomService) CreateRoom(ctx context.Context, req room.CreateRoomReq) error {
	id, err := r.storage.CreateRoom(ctx, req)
	if err != nil {
		return err
	}
	room := &domain.Room{
		ID:       id,
		Name:     req.Name,
		Password: req.Password,
	}
	return r.pusher.Produce(room, roomTopic, roomKey)
}
