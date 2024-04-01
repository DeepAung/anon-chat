package main

import (
	"net/http"

	"github.com/DeepAung/anon-chat/server/handlers"
	"github.com/DeepAung/anon-chat/server/hub"
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

func (r *router) UsersRouter(mux *http.ServeMux) {
	handler := handlers.NewUsersHandler()
	mux.HandleFunc("POST /api/users/login", handler.Login)
	mux.HandleFunc("GET /api/users/logout", handler.Logout)
}

func (r *router) PagesRouter(mux *http.ServeMux) {
	handler := handlers.NewPagesHandler()

	mux.HandleFunc("GET /index", handler.Index)
	mux.HandleFunc("GET /login", handler.Login)
}

func (r *router) TestRouter(mux *http.ServeMux, hub *hub.Hub) {
	mux.HandleFunc("GET /api/test/rooms", func(w http.ResponseWriter, r *http.Request) {
		w.Write(hub.RoomsMarshalled())
	})
}
