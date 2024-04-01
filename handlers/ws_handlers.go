package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DeepAung/anon-chat/server/hub"
	"github.com/DeepAung/anon-chat/server/types"
	"github.com/DeepAung/anon-chat/server/utils"
	"golang.org/x/net/websocket"
)

type wsHandler struct {
	hub *hub.Hub
}

func NewWsHandler(hub *hub.Hub) *wsHandler {
	return &wsHandler{
		hub: hub,
	}
}

func (h *wsHandler) Connect(w http.ResponseWriter, r *http.Request) {
	username := utils.GetCookieValue(r, "username")
	if username == "" {
		http.Error(w, "cookie not found", http.StatusBadRequest)
		return
	}

	roomId := r.PathValue("roomId")

	fmt.Println("user: ", username, " | roomId: ", roomId)

	websocket.Handler(func(ws *websocket.Conn) {
		err := h.hub.Connect(ws, username, roomId)
		if err != nil {
			wsMsg := types.WsMessage{
				Type:    types.IsError,
				Content: err.Error(),
			}
			res, _ := json.Marshal(&wsMsg)

			ws.Write(res)
			ws.Close()
			return
		}

		h.hub.Listen(ws, roomId)
	}).ServeHTTP(w, r)
}

func (h *wsHandler) CreateAndConnect(w http.ResponseWriter, r *http.Request) {
	username := utils.GetCookieValue(r, "username")
	if username == "" {
		http.Error(w, "cookie not found", http.StatusBadRequest)
		return
	}

	roomName := r.PathValue("roomName")

	fmt.Println("username: ", username, " | roomName: ", roomName)

	websocket.Handler(func(ws *websocket.Conn) {
		roomId := h.hub.CreateAndConnect(ws, username, roomName)
		h.hub.Listen(ws, roomId)
	}).ServeHTTP(w, r)
}
