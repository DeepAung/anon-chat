package types

import (
	"golang.org/x/net/websocket"
)

type Room struct {
	Id         string                   `json:"id"`
	Name       string                   `json:"name"`
	Users      map[*websocket.Conn]User `json:"users"`
	MsgHistory *History                 `json:"messages"`
}

func NewRoom(id, name string, historyLen int) Room {
	return Room{
		Id:         id,
		Name:       name,
		Users:      make(map[*websocket.Conn]User),
		MsgHistory: NewHistory(historyLen),
	}
}
