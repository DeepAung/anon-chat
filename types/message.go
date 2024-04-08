package types

import (
	"time"

	"github.com/google/uuid"
)

type ReqMessage struct {
	Type    MessageType `json:"type"`
	Content string      `json:"content"`
}

type MessageType string

const (
	UserType       MessageType = "user"
	SystemType     MessageType = "system"
	DisconnectType MessageType = "disconnect"
)

type ResMessage struct {
	Id      string      `json:"id"`
	Type    MessageType `json:"type"`
	User    User        `json:"user"`
	Content string      `json:"content"`
	Time    time.Time   `json:"time"`
}

func NewResMessage(msgType MessageType, user User, content string) ResMessage {
	return ResMessage{
		Id:      uuid.NewString(),
		Type:    msgType,
		User:    user,
		Content: content,
		Time:    time.Now(),
	}
}
