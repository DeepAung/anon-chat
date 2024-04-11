package server

import (
	"log"
	"net/http"

	"github.com/DeepAung/anon-chat/pkg/config"
	"github.com/DeepAung/anon-chat/pkg/hub"
	"github.com/DeepAung/anon-chat/pkg/router"
)

type Server struct {
	r   *router.Router
	cfg *config.Config
	hub *hub.Hub
}

func NewServer(r *router.Router, cfg *config.Config, hub *hub.Hub) *Server {
	return &Server{
		r:   r,
		cfg: cfg,
		hub: hub,
	}
}

func (s *Server) Start() {
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
