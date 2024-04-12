package router

import (
	"net/http"

	"github.com/DeepAung/anon-chat/handlers"
	"github.com/DeepAung/anon-chat/pkg/hub"
)

type Router struct{}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) WsRouter(mux *http.ServeMux, hub *hub.Hub) {
	handler := handlers.NewWsHandler(hub)

	mux.HandleFunc("/ws/connect", handler.Connect)
}

func (r *Router) RoomsRouter(mux *http.ServeMux, hub *hub.Hub) {
	handler := handlers.NewRoomsHandler(hub)

	mux.HandleFunc("POST /api/rooms/create-and-connect", handler.CreateAndConnect)
	mux.HandleFunc("POST /api/rooms/connect", handler.Connect)
}

func (r *Router) TestRouter(mux *http.ServeMux, hub *hub.Hub) {
	mux.HandleFunc("GET /api/test/rooms", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(hub.JSONRooms())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func (r *Router) PagesRouter(mux *http.ServeMux, hub *hub.Hub) {
	handler := handlers.NewPagesHandler(hub)

	mux.HandleFunc("GET /chat", handler.Chat)
	mux.HandleFunc("/", handler.Index)
}
