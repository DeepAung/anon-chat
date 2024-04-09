package main

import (
	"log"
	"net/http"

	"github.com/DeepAung/anon-chat/hub"
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
	s.r.RoomsRouter(mux, s.hub)
	s.r.TestRouter(mux, s.hub)
	s.r.PagesRouter(mux, s.hub)

	// TODO: change to ":3000" in prod (I add localhost bc i don't wanna fight windows defender firewall everytime code is updated)
	log.Fatal(http.ListenAndServe("127.0.0.1:3000", mux))
}
