package main

import (
	"github.com/DeepAung/anon-chat/server/hub"
)

func main() {
	startServer()
}

func startServer() {
	hub := hub.NewHub()

	router := NewRouter()

	server := NewServer(router, hub)
	server.Start()
}
