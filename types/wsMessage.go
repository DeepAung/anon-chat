package types

type WsMessage struct {
	Type    WsMessageType `json:"type"`
	Content string        `json:"content"`
}

type WsMessageType string

const (
	IsMessage WsMessageType = "message"
	IsError   WsMessageType = "error"
)
