package hub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/DeepAung/anon-chat/types"
	"github.com/DeepAung/anon-chat/views"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

var RoomIdNotFoundErr error = errors.New("room id not found")

type Hub struct {
	mu    sync.Mutex
	rooms map[string]*types.Room
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]*types.Room),
	}
}

func (h *Hub) RoomsMarshalled() []byte {
	tmp := make(map[string]struct {
		Id    string       `json:"id"`
		Name  string       `json:"name"`
		Users []types.User `json:"users"`
	})

	for key, value := range h.rooms {
		users := []types.User{}
		for _, user := range value.Users {
			users = append(users, user)
		}

		tmp[key] = struct {
			Id    string       `json:"id"`
			Name  string       `json:"name"`
			Users []types.User `json:"users"`
		}{
			Id:    value.Id,
			Name:  value.Name,
			Users: users,
		}
	}

	res, _ := json.Marshal(tmp)
	return res
}

func (h *Hub) Create(roomName string) string {
	roomId := uuid.NewString()
	newRoom := &types.Room{
		Id:    roomId,
		Name:  roomName,
		Users: make(map[*websocket.Conn]types.User),
	}

	h.mu.Lock()
	h.rooms[roomId] = newRoom
	h.mu.Unlock()

	return roomId
}

func (h *Hub) Connect(conn *websocket.Conn, username string, roomId string) error {
	h.mu.Lock()
	room, ok := h.rooms[roomId]
	if !ok {
		h.mu.Unlock()
		return RoomIdNotFoundErr
	}

	room.Users[conn] = types.User{
		Id:       uuid.NewString(),
		Username: username,
	}

	h.mu.Unlock()
	return nil
}

func (h *Hub) Disconnect(conn *websocket.Conn, roomId string) error {
	h.mu.Lock()
	room, ok := h.rooms[roomId]
	if !ok {
		h.mu.Unlock()
		return RoomIdNotFoundErr
	}

	delete(room.Users, conn)
	if len(room.Users) == 0 {
		delete(h.rooms, roomId)
	}

	h.mu.Unlock()
	return nil
}

func (h *Hub) Broadcast(msg types.ResMessage, roomId string) error {
	h.mu.Lock()
	for conn := range h.rooms[roomId].Users {
		if err := views.Message(msg).Render(context.Background(), conn); err != nil {
			h.mu.Unlock()
			return err
		}
	}

	h.mu.Unlock()
	return nil
}

func (h *Hub) ConnectAndListen(conn *websocket.Conn, username, roomId string) {
	defer conn.Close()

	if err := h.Connect(conn, username, roomId); err != nil {
		views.ErrorMessage(err.Error()).Render(context.Background(), conn)
		return
	}

	resMsg := types.NewSystemMessage(username + " is joined")
	h.Broadcast(resMsg, roomId)

	h.mu.Lock()
	user := h.rooms[roomId].Users[conn]
	h.mu.Unlock()
	var reqMsg types.ReqMessage

	for {
		if err := websocket.JSON.Receive(conn, &reqMsg); err != nil {
			if err == io.EOF {
				continue
			}

			fmt.Println("disconnect bz error: ", err)

			views.ErrorMessage(err.Error()).Render(context.Background(), conn)

			resMsg := types.NewSystemMessage(username + " is leaved")
			h.Broadcast(resMsg, roomId)
			h.Disconnect(conn, roomId)
			return
		}

		if reqMsg.Type == types.DisconnectType {
			resMsg := types.NewSystemMessage(username + " is leaved")
			h.Broadcast(resMsg, roomId)
			h.Disconnect(conn, roomId)
			return
		}

		resMsg := types.NewResMessage(reqMsg.Type, user, reqMsg.Content)
		h.Broadcast(resMsg, roomId)
	}
}
