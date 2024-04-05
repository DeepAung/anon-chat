package types

import (
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type Room struct {
	Id    string                   `json:"id"`
	Name  string                   `json:"name"`
	Users map[*websocket.Conn]User `json:"users"`
	// Messages []Message `json:"messages"`
}

func NewRoom(name string) Room {
	return Room{
		Id:    uuid.NewString(),
		Name:  name,
		Users: make(map[*websocket.Conn]User),
	}
}
