package clientstor

import (
	"app-websocket/internal/domain"
	"app-websocket/internal/storage/connector"
	"context"
	"errors"
	"fmt"
)

type ClientRepository struct {
	db *connector.Database
}

func NewClientRepository(db *connector.Database) *ClientRepository {
	return &ClientRepository{
		db: db,
	}
}

func (r *ClientRepository) AddClient(ctx context.Context, req *domain.Client) error {
	op := "repository.AddClient"

	clQuery := "INSERT INTO clients (id, username, room_id) SELECT $1, $2, $3 WHERE NOT EXISTS (SELECT 1 FROM clients WHERE id = $1 AND room_id = $3)"
	row, err := r.db.GetDB().ExecContext(ctx, clQuery, req.ClientID, req.Username, req.RoomID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	res, err := row.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if res > 1 {
		return fmt.Errorf("%s: %w", op, errors.New("error while adding client in db"))
	}
	return nil
}
