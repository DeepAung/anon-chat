package hub

import (
	"fmt"
	"io"

	"github.com/DeepAung/anon-chat/server/types"
	"golang.org/x/net/websocket"
)

type connReq struct {
	conn   *websocket.Conn
	roomId string
}

type broadcastReq struct {
	msg    types.Message
	roomId string
}

type Hub struct {
	rooms          map[string]map[*websocket.Conn]struct{}
	connectChan    chan connReq
	disconnectChan chan connReq
	broadcastChan  chan broadcastReq
}

func NewHub() *Hub {
	return &Hub{
		make(map[string]map[*websocket.Conn]struct{}),
		make(chan connReq),
		make(chan connReq),
		make(chan broadcastReq),
	}
}

func (h *Hub) RunLoop() {
	for {
		select {
		case req := <-h.connectChan:
			h.connect(req.conn, req.roomId)
		case req := <-h.disconnectChan:
			h.disconnect(req.conn, req.roomId)
		case req := <-h.broadcastChan:
			h.broadcast(req.msg, req.roomId)
		}
	}
}

func (h *Hub) Connect(conn *websocket.Conn, roomId string) {
	h.connectChan <- connReq{conn: conn, roomId: roomId}
}

func (h *Hub) Disconnect(conn *websocket.Conn, roomId string) {
	h.disconnectChan <- connReq{conn: conn, roomId: roomId}
}

func (h *Hub) Broadcast(msg types.Message, roomId string) {
	h.broadcastChan <- broadcastReq{msg: msg, roomId: roomId}
}

func (h *Hub) connect(conn *websocket.Conn, roomId string) {
	if _, ok := h.rooms[roomId]; !ok {
		h.rooms[roomId] = make(map[*websocket.Conn]struct{})
	}
	h.rooms[roomId][conn] = struct{}{}

	fmt.Printf("[%s|%s] connect to [%s]\n", conn.RemoteAddr().String(), conn.LocalAddr().String(), roomId)
}

func (h *Hub) disconnect(conn *websocket.Conn, roomId string) {
	delete(h.rooms[roomId], conn)
	if len(h.rooms[roomId]) == 0 {
		delete(h.rooms, roomId)
	}

	fmt.Printf("[%s|%s] disconnect to [%s]\n", conn.RemoteAddr().String(), conn.LocalAddr().String(), roomId)
}

func (h *Hub) broadcast(msg types.Message, roomId string) {
	for conn := range h.rooms[roomId] {
		if err := websocket.JSON.Send(conn, msg); err != nil {
			fmt.Printf("broadcast message failed: %v\n", err)
			return
		}
	}

	fmt.Printf("[%s] is broadcasted to [%s]\n", msg, roomId)
}

func (h *Hub) ConnectAndListen(ws *websocket.Conn, roomId string) {
	h.Connect(ws, roomId)

	var msg types.Message
	for {
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			if err == io.EOF {
				continue
			}

			h.Disconnect(ws, roomId)
			fmt.Println("error: ", err)
			continue
		}

		h.Broadcast(msg, roomId)
	}
}
