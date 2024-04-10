package main

import (
	"log"
	"net/http"

	"github.com/DeepAung/anon-chat/hub"
)

type server struct {
	r   *router
	cfg *config
	hub *hub.Hub
}

func NewServer(r *router, cfg *config, hub *hub.Hub) *server {
	return &server{
		r:   r,
		cfg: cfg,
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

	port := ":3000"
	if !s.cfg.IsProd {
		port = "127.0.0.1:3000"
	}

	log.Fatal(http.ListenAndServe(port, mux))
}
