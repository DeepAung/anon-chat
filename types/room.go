package types

import "golang.org/x/net/websocket"

type Room struct {
	Id    string                   `json:"id"`
	Name  string                   `json:"name"`
	Users map[*websocket.Conn]User `json:"users"`
	// Messages []Message `json:"messages"`
}
