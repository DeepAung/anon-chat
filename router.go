package main

import (
	"net/http"

	"github.com/DeepAung/anon-chat/handlers"
	"github.com/DeepAung/anon-chat/hub"
	"github.com/DeepAung/anon-chat/services"
)

type router struct{}

func NewRouter() *router {
	return &router{}
}

func (r *router) WsRouter(mux *http.ServeMux, hub *hub.Hub) {
	handler := handlers.NewWsHandler(hub)

	mux.HandleFunc("/ws/connect/{roomId}", handler.Connect)
	mux.HandleFunc("/ws/create-and-connect/{roomName}", handler.CreateAndConnect)
}

func (r *router) RoomsRouter(mux *http.ServeMux, hub *hub.Hub) {
	usersSvc := services.NewUsersService()
	handler := handlers.NewRoomsHandler(hub, usersSvc)

	mux.HandleFunc("POST /api/rooms/create-connect", handler.CreateConnect)
	mux.HandleFunc("POST /api/rooms/disconnect", handler.Disconnect)
}

func (r *router) TestRouter(mux *http.ServeMux, hub *hub.Hub) {
	mux.HandleFunc("GET /api/test/rooms", func(w http.ResponseWriter, r *http.Request) {
		w.Write(hub.RoomsMarshalled())
	})
}

func (r *router) PagesRouter(mux *http.ServeMux) {
	handler := handlers.NewPagesHandler()

	mux.HandleFunc("GET /chat", handler.Chat)
	mux.HandleFunc("/", handler.Index)
}
