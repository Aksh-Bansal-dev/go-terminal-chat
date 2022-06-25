package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Aksh-Bansal-dev/go-terminal-chat/internal/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	hub := websocket.NewHub()
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	})
	log.Println(fmt.Sprintf("Server started at %s", *addr))
	log.Fatal(http.ListenAndServe(*addr, nil))
}
