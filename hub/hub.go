package hub

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/DeepAung/anon-chat/types"
	"github.com/DeepAung/anon-chat/views"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

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
	if room, ok := h.rooms[roomId]; ok {
		room.Users[conn] = types.User{
			Id:       uuid.NewString(),
			Username: username,
		}

		h.mu.Unlock()

		content := username + " is joined"
		h.Broadcast(types.NewResMessage(types.SystemType, types.SystemUser, content), roomId)
		// fmt.Println(username, "connect to ", roomId)
		return nil
	}

	h.mu.Unlock()

	// fmt.Println("error: room id not found")
	return fmt.Errorf("room id not found")
}

func (h *Hub) Disconnect(conn *websocket.Conn, roomId string) error {
	h.mu.Lock()
	if room, ok := h.rooms[roomId]; ok {
		delete(room.Users, conn)
		if len(room.Users) == 0 {
			delete(h.rooms, roomId)
		}

		h.mu.Unlock()

		content := room.Users[conn].Username + " is leaved"
		h.Broadcast(types.NewResMessage(types.SystemType, types.SystemUser, content), roomId)
		// fmt.Println("disconnect to ", roomId)
		return nil
	}

	h.mu.Unlock()
	return fmt.Errorf("room id not found")
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
	// fmt.Println("roomId: ", roomId, " | msg: ", msg, "| is broadcasted")
	return nil
}

func (h *Hub) Listen(conn *websocket.Conn, roomId string) {
	user := h.rooms[roomId].Users[conn]
	reqMsg := new(types.ReqMessage)

	// fmt.Println("listening to ", roomId)
	for {
		if err := websocket.JSON.Receive(conn, reqMsg); err != nil {
			if err == io.EOF {
				continue
			}

			h.Disconnect(conn, roomId)
			// fmt.Println("error: ", err)
			continue
		}

		resMsg := types.NewResMessage(reqMsg.Type, user, reqMsg.Content)
		h.Broadcast(resMsg, roomId)
	}
}
