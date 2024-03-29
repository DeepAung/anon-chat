package types

import (
	"time"
)

type Message struct {
	Username string    `json:"username"`
	Content  string    `json:"content"`
	Time     time.Time `json:"time"`
}

type Room struct {
	Name string
	Uuid string
}
