package room

import (
	"app-websocket/internal/domain"

	"app-websocket/internal/domain/room"
	"context"
)

type RoomStorage interface {
	CreateRoom(ctx context.Context, req room.CreateRoomReq) (int64, error)
	GetRooms(ctx context.Context) ([]domain.Room, error)
}

type RoomService struct {
	r_storage RoomStorage
}

func NewRoomService(rstorage RoomStorage) *RoomService {
	return &RoomService{

		r_storage: rstorage,
	}
}

func (r *RoomService) CreateRoom(ctx context.Context, req room.CreateRoomReq) (int64, error) {
	id, err := r.r_storage.CreateRoom(ctx, req)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *RoomService) GetRooms(ctx context.Context) ([]domain.Room, error) {
	data, err := r.r_storage.GetRooms(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}
