package hub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/DeepAung/anon-chat/pkg/config"
	"github.com/DeepAung/anon-chat/types"
	"github.com/DeepAung/anon-chat/views"
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

var RoomIdNotFoundErr error = errors.New("room id not found")

type Hub struct {
	cfg *config.Config

	mu    sync.Mutex
	rooms map[string]*types.Room
}

func NewHub(cfg *config.Config) *Hub {
	return &Hub{
		cfg:   cfg,
		rooms: make(map[string]*types.Room),
	}
}

func (h *Hub) JSONRooms() []byte {
	type tmpStruct struct {
		Id       string             `json:"id"`
		Name     string             `json:"name"`
		Users    []types.User       `json:"users"`
		Messages []types.ResMessage `json:"messages"`
	}

	tmp := make([]tmpStruct, 0)

	h.mu.Lock()
	for _, value := range h.rooms {
		users := make([]types.User, 0)
		for _, user := range value.Users {
			users = append(users, user)
		}

		messages := make([]types.ResMessage, 0)
		iter := value.MsgHistory.Iter()
		for iter.Next() {
			messages = append(messages, iter.Get())
		}

		tmp = append(tmp, tmpStruct{
			Id:       value.Id,
			Name:     value.Name,
			Users:    users,
			Messages: messages,
		})
	}
	h.mu.Unlock()

	res, _ := json.Marshal(tmp)
	return res
}

func (h *Hub) GetRoom(roomId string) (types.Room, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, ok := h.rooms[roomId]
	if !ok {
		return types.Room{}, false
	}
	return *room, true
}

func (h *Hub) CreateRoom(roomName string) string {
	roomId := uuid.NewString()
	newRoom := types.NewRoom(roomId, roomName, h.cfg.HistoryLength)

	h.mu.Lock()
	h.rooms[roomId] = &newRoom
	h.mu.Unlock()

	return roomId
}

func (h *Hub) ConnectAndListen(conn *websocket.Conn, username, roomId string) {
	defer conn.Close()

	if err := h.connect(conn, username, roomId); err != nil {
		views.ErrorBody(err.Error()).Render(context.Background(), conn)
		return
	}

	h.mu.Lock()
	user := h.rooms[roomId].Users[conn]
	h.mu.Unlock()

	resMsg := types.NewSystemMessage(username + " is joined")
	h.broadcastMessage(resMsg, roomId)
	h.broadcastMemberJoin(user, roomId)

	var reqMsg types.ReqMessage
	for {
		if err := websocket.JSON.Receive(conn, &reqMsg); err != nil {
			if err == io.EOF {
				continue
			}

			fmt.Println("disconnect bz error: ", err)
			views.ErrorBody(err.Error()).Render(context.Background(), conn)

			resMsg := types.NewSystemMessage(username + " is leaved")
			h.broadcastMessage(resMsg, roomId)
			h.broadcastMemberLeave(user, roomId)

			h.disconnect(conn, roomId)
			return
		}

		if reqMsg.Type == types.DisconnectType {
			resMsg := types.NewSystemMessage(username + " is leaved")
			h.broadcastMessage(resMsg, roomId)
			h.broadcastMemberLeave(user, roomId)

			h.disconnect(conn, roomId)
			return
		}

		resMsg := types.NewResMessage(reqMsg.Type, user, reqMsg.Content)
		h.broadcastMessage(resMsg, roomId)
	}
}

func (h *Hub) connect(conn *websocket.Conn, username string, roomId string) error {
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

func (h *Hub) disconnect(conn *websocket.Conn, roomId string) error {
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

func (h *Hub) broadcastMessage(msg types.ResMessage, roomId string) error {
	h.mu.Lock()
	h.rooms[roomId].MsgHistory.Push(msg)
	h.mu.Unlock()

	return h.broadcast(views.Message(msg), roomId)
}

func (h *Hub) broadcastMemberJoin(user types.User, roomId string) error {
	return h.broadcast(views.MemberJoin(user), roomId)
}

func (h *Hub) broadcastMemberLeave(user types.User, roomId string) error {
	return h.broadcast(views.MemberLeave(user), roomId)
}

func (h *Hub) broadcast(component templ.Component, roomId string) error {
	h.mu.Lock()
	for conn := range h.rooms[roomId].Users {
		if err := component.Render(context.Background(), conn); err != nil {
			h.mu.Unlock()
			return err
		}
	}

	h.mu.Unlock()
	return nil
}
