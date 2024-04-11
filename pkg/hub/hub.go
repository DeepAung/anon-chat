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
	roomsMu sync.Mutex
	rooms   map[string]*types.Room

	historyMu sync.Mutex
	history   *types.History
}

func NewHub(cfg *config.Config) *Hub {
	return &Hub{
		rooms:   make(map[string]*types.Room),
		history: types.NewHistory(int(cfg.HistoryLength)),
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

func (h *Hub) GetRoom(roomId string) (types.Room, bool) {
	h.roomsMu.Lock()
	defer h.roomsMu.Unlock()

	room, ok := h.rooms[roomId]
	if !ok {
		return types.Room{}, false
	}
	return *room, true
}

func (h *Hub) CreateRoom(roomName string) string {
	roomId := uuid.NewString()
	newRoom := &types.Room{
		Id:         roomId,
		Name:       roomName,
		Users:      make(map[*websocket.Conn]types.User),
		MsgHistory: h.history,
	}

	h.roomsMu.Lock()
	h.rooms[roomId] = newRoom
	h.roomsMu.Unlock()

	return roomId
}

func (h *Hub) ConnectAndListen(conn *websocket.Conn, username, roomId string) {
	defer conn.Close()

	if err := h.connect(conn, username, roomId); err != nil {
		views.ErrorBody(err.Error()).Render(context.Background(), conn)
		return
	}

	h.roomsMu.Lock()
	user := h.rooms[roomId].Users[conn]
	h.roomsMu.Unlock()

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
	h.roomsMu.Lock()
	room, ok := h.rooms[roomId]
	if !ok {
		h.roomsMu.Unlock()
		return RoomIdNotFoundErr
	}

	room.Users[conn] = types.User{
		Id:       uuid.NewString(),
		Username: username,
	}

	h.roomsMu.Unlock()
	return nil
}

func (h *Hub) disconnect(conn *websocket.Conn, roomId string) error {
	h.roomsMu.Lock()
	room, ok := h.rooms[roomId]
	if !ok {
		h.roomsMu.Unlock()
		return RoomIdNotFoundErr
	}

	delete(room.Users, conn)
	if len(room.Users) == 0 {
		delete(h.rooms, roomId)
	}

	h.roomsMu.Unlock()
	return nil
}

func (h *Hub) broadcast(component templ.Component, roomId string) error {
	h.roomsMu.Lock()
	for conn := range h.rooms[roomId].Users {
		if err := component.Render(context.Background(), conn); err != nil {
			h.roomsMu.Unlock()
			return err
		}
	}

	h.roomsMu.Unlock()
	return nil
}

func (h *Hub) broadcastMessage(msg types.ResMessage, roomId string) error {
	return h.broadcast(views.Message(msg), roomId)
}

func (h *Hub) broadcastMemberJoin(user types.User, roomId string) error {
	return h.broadcast(views.MemberJoin(user), roomId)
}

func (h *Hub) broadcastMemberLeave(user types.User, roomId string) error {
	return h.broadcast(views.MemberLeave(user), roomId)
}
