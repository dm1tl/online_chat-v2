package room

import (
	"app-websocket/internal/domain/room"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type RoomRepository struct {
	db *sqlx.DB
}

func NewRoomRepository(db *sqlx.DB) *RoomRepository {
	return &RoomRepository{
		db: db,
	}
}

func (r *RoomRepository) CreateRoom(ctx context.Context, req room.CreateRoomReq) (int64, error) {
	op := "repository.CreateRoom"
	query := "INSERT INTO rooms (name, password) VALUES ($1, $2) RETURNING id"
	var id int64

	if err := r.db.QueryRowContext(ctx, query, req.Name, req.Password).Scan(&id); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}
