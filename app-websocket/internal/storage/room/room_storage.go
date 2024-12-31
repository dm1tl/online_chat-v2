package roomstor

import (
	"app-websocket/internal/domain"
	"app-websocket/internal/domain/room"
	"app-websocket/internal/storage/connector"
	"context"
	"fmt"
)

type RoomRepository struct {
	db *connector.Database
}

func NewRoomRepository(db *connector.Database) *RoomRepository {
	return &RoomRepository{
		db: db,
	}
}

func (r *RoomRepository) CreateRoom(ctx context.Context, req room.CreateRoomReq) (int64, error) {
	op := "repository.CreateRoom"
	query := "INSERT INTO rooms (name, password) VALUES ($1, $2) RETURNING id"
	var id int64

	if err := r.db.GetDB().QueryRowContext(ctx, query, req.Name, req.Password).Scan(&id); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (r *RoomRepository) GetRooms(ctx context.Context) ([]domain.Room, error) {
	op := "repository.GetAllRooms"
	var output []domain.Room
	query := "SELECT id, name, password FROM rooms"
	res, err := r.db.GetDB().QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer res.Close()
	for res.Next() {
		var room domain.Room
		if err := res.Scan(&room.ID, &room.Name, &room.Password); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		output = append(output, room)
	}
	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return output, nil
}
