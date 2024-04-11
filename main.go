package main

import (
	"github.com/DeepAung/anon-chat/pkg/config"
	"github.com/DeepAung/anon-chat/pkg/hub"
	"github.com/DeepAung/anon-chat/pkg/router"
	"github.com/DeepAung/anon-chat/pkg/server"
)

func main() {
	startServer()
}

func startServer() {
	cfg := config.LoadConfig()
	router := router.NewRouter()
	hub := hub.NewHub(cfg)

	server := server.NewServer(router, cfg, hub)
	server.Start()
}
