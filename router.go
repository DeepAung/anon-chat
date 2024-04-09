package main

import (
	"net/http"

	"github.com/DeepAung/anon-chat/handlers"
	"github.com/DeepAung/anon-chat/hub"
)

type router struct{}

func NewRouter() *router {
	return &router{}
}

func (r *router) WsRouter(mux *http.ServeMux, hub *hub.Hub) {
	handler := handlers.NewWsHandler(hub)

	mux.HandleFunc("/ws/connect", handler.Connect)
}

func (r *router) RoomsRouter(mux *http.ServeMux, hub *hub.Hub) {
	handler := handlers.NewRoomsHandler(hub)

	mux.HandleFunc("POST /api/rooms/create-and-connect", handler.CreateAndConnect)
	mux.HandleFunc("POST /api/rooms/connect", handler.Connect)
}

func (r *router) TestRouter(mux *http.ServeMux, hub *hub.Hub) {
	mux.HandleFunc("GET /api/test/rooms", func(w http.ResponseWriter, r *http.Request) {
		w.Write(hub.RoomsMarshalled())
	})
}

func (r *router) PagesRouter(mux *http.ServeMux, hub *hub.Hub) {
	handler := handlers.NewPagesHandler(hub)

	mux.HandleFunc("GET /chat", handler.Chat)
	mux.HandleFunc("/", handler.Index)
}
