package types

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Id      string    `json:"id"`
	User    User      `json:"user"`
	Content string    `json:"content"`
	Time    time.Time `json:"time"`
}

func NewMessage(user User, content string) Message {
	return Message{
		Id:      uuid.NewString(),
		User:    user,
		Content: content,
		Time:    time.Now(),
	}
}
