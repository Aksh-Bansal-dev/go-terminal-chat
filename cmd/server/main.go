package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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
			return
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
		jsonRes, err := json.Marshal(onlineUsers)
		if err != nil {
			log.Println(err)
		}
		w.Write(jsonRes)
	})
	// new htttp endpoint to check if username is present in clients
	http.HandleFunc("/valid-username", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Cannot parse request body", http.StatusBadRequest)
		}
		var data map[string]string
		err = json.Unmarshal(body, &data)
		if err != nil {
			http.Error(w, "Cannot parse request body", http.StatusBadRequest)
		}
		for client := range hub.Clients {
			if client.Username == data["username"] {
				res := map[string]bool{"valid": false}
				jsonRes, _ := json.Marshal(res)
				w.Write([]byte(jsonRes))
				return
			}
		}
		res := map[string]bool{"valid": true}
		jsonRes, _ := json.Marshal(res)
		w.Write([]byte(jsonRes))
	})
	log.Println(fmt.Sprintf("Server started at %s", *addr))
	log.Fatal(http.ListenAndServe(*addr, nil))
}
