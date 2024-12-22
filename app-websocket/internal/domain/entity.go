package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

type Message struct {
	Content     string
	Username    string
	TimeCreated time.Time
	RoomID      string
	UserID      string
}

type MessageHandler func(msg Message) error

func NewMessage(value []byte) (*Message, error) {
	op := "domain.NewMessage"
	var msg Message
	if err := json.Unmarshal(value, &msg); err != nil {
		return nil, fmt.Errorf("op: %s: %w", op, err)
	}
	return &msg, nil
}
