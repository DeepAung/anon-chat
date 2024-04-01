package hub

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/DeepAung/anon-chat/server/types"
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

func (h *Hub) CreateAndConnect(conn *websocket.Conn, username, roomName string) string {
	roomId := uuid.NewString()
	newRoom := &types.Room{
		Id:    roomId,
		Name:  roomName,
		Users: make(map[*websocket.Conn]types.User),
	}
	newRoom.Users[conn] = types.User{
		Id:       uuid.NewString(),
		Username: username,
	}

	h.mu.Lock()
	h.rooms[roomId] = newRoom
	h.mu.Unlock()

	fmt.Println(username, "create and connect to ", roomName)
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
		fmt.Println(username, "connect to ", roomId)
		return nil
	}

	h.mu.Unlock()
	fmt.Println("error: room id not found")
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
		return nil
	}

	h.mu.Unlock()
	fmt.Println("disconnect to ", roomId)
	return fmt.Errorf("room id not found")
}

func (h *Hub) Broadcast(msg types.Message, roomId string) error {
	h.mu.Lock()
	for conn := range h.rooms[roomId].Users {
		if err := websocket.JSON.Send(conn, msg); err != nil {
			h.mu.Unlock()
			return err
		}
	}

	h.mu.Unlock()
	fmt.Println("in ", roomId, " |", msg, "| is broadcasted")
	return nil
}

func (h *Hub) Listen(conn *websocket.Conn, roomId string) {
	user := h.rooms[roomId].Users[conn]
	var content string

	fmt.Println("listening to ", roomId)
	for {
		if err := websocket.JSON.Receive(conn, &content); err != nil {
			if err == io.EOF {
				continue
			}

			h.Disconnect(conn, roomId)
			fmt.Println("error: ", err)
			continue
		}

		msg := types.NewMessage(user, content)
		h.Broadcast(msg, roomId)
	}
}
