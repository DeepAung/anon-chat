package main

import (
	"github.com/DeepAung/anon-chat/hub"
)

func main() {
	startServer()
}

func startServer() {
	router := NewRouter()
	cfg := LoadConfig()
	hub := hub.NewHub()

	server := NewServer(router, cfg, hub)
	server.Start()
}
