package handlers

import (
	"net/http"

	"github.com/DeepAung/anon-chat/server/hub"
	"golang.org/x/net/websocket"
)

func WsHandler(mux *http.ServeMux, h *hub.Hub) {
	mux.HandleFunc("/ws/{roomId}", func(w http.ResponseWriter, r *http.Request) {
		roomId := r.PathValue("roomId")

		websocket.Handler(func(ws *websocket.Conn) {
			h.ConnectAndListen(ws, roomId)
		}).ServeHTTP(w, r)

	})
}
