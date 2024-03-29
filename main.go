package main

import (
	"log"
	"net/http"

	"github.com/DeepAung/anon-chat/server/handlers"
	"github.com/DeepAung/anon-chat/server/hub"
)

func main() {
	startServer()
}

func startServer() {
	h := hub.NewHub()
	go h.RunLoop()

	mux := http.NewServeMux()

	handlers.WsHandler(mux, h)
	handlers.PagesHandler(mux)

	log.Fatal(http.ListenAndServe(":3000", mux))
}
