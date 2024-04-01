package main

import (
	"log"
	"net/http"

	"github.com/DeepAung/anon-chat/server/hub"
)

type server struct {
	r   *router
	hub *hub.Hub
}

func NewServer(r *router, hub *hub.Hub) *server {
	return &server{
		r:   r,
		hub: hub,
	}
}

func (s *server) Start() {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	s.r.WsRouter(mux, s.hub)
	s.r.UsersRouter(mux)
	s.r.PagesRouter(mux)
	s.r.TestRouter(mux, s.hub)

	log.Fatal(http.ListenAndServe(":3000", mux))
}
