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

type Message struct {
	Content     string    `json:"content"`
	Username    string    `json:"username"`
	TimeCreated time.Time `json:"time_created"`
	RoomID      string    `json:"room_id"`
	UserID      string    `json:"user_id"`
}

func (m *Message) GetEvent() interface{} {
	return m
}

type EventHandler func(msg Event) error

func NewEvent(value []byte) (Event, error) {
	op := "domain.NewEvent"
	var msg Message
	if err := json.Unmarshal(value, &msg); err == nil {
		return &msg, nil
	}

	var room Room
	if err := json.Unmarshal(value, &room); err == nil {
		return &room, nil
	}
	return nil, fmt.Errorf("op: %s: %w", op, errCouldNotSerialize)
}

func NewMessage(value []byte) (*Message, error) {
	op := "domain.NewMessage"
	var msg Message
	if err := json.Unmarshal(value, &msg); err != nil {
		return nil, fmt.Errorf("op: %s: %w", op, err)
	}
	return &msg, nil
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Room struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (r *Room) GetEvent() interface{} {
	return r
}
