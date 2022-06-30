package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Aksh-Bansal-dev/go-terminal-chat/internal/user"
	"github.com/Aksh-Bansal-dev/go-terminal-chat/internal/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile)
	hub := websocket.NewHub()
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	})
	http.HandleFunc("/online-users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		}
		w.Header().Set("Content-Type", "application/json")
		onlineUsers := []user.OnlineUser{}
		for client := range hub.Clients {
			onlineUsers = append(
				onlineUsers,
				user.OnlineUser{
					Username: (*client).Username,
					Color:    (*client).Color,
				},
			)
		}
		rawData, err := json.Marshal(onlineUsers)
		if err != nil {
			log.Println(err)
		}
		w.Write(rawData)
	})
	log.Println(fmt.Sprintf("Server started at %s", *addr))
	log.Fatal(http.ListenAndServe(*addr, nil))
}
