package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var errCouldNotSerialize = errors.New("could not serialize event")

type Event interface {
	GetEvent() interface{}
}

type EventHandler func(Event) error

type MessageHandler func(msg Message) error

type Message struct {
	Content     string    `json:"content"`
	Username    string    `json:"username"`
	TimeCreated time.Time `json:"time_created"`
	RoomID      int64     `json:"room_id"`
	UserID      int64     `json:"user_id"`
}

func (m *Message) GetEvent() interface{} {
	return m
}

type Room struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func NewEvent(value []byte) (Event, error) {
	op := "domain.NewEvent"
	var msg Message
	if err := json.Unmarshal(value, &msg); err == nil {
		return &msg, nil
	}
	var client Client
	if err := json.Unmarshal(value, &client); err == nil {
		return &client, err
	}

	return nil, fmt.Errorf("op: %s: %w", op, errCouldNotSerialize)
}

type Client struct {
	RoomID   int64  `json:"room_id"`
	ClientID int64  `json:"client_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (m *Client) GetEvent() interface{} {
	return m
}
