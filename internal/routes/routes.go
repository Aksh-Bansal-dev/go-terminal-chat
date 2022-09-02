package routes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Aksh-Bansal-dev/go-terminal-chat/internal/database"
	"github.com/Aksh-Bansal-dev/go-terminal-chat/internal/user"
	"github.com/Aksh-Bansal-dev/go-terminal-chat/internal/websocket"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func OnlineUserHandler(w http.ResponseWriter, r *http.Request, hub *websocket.Hub) {
	if r.Method != "GET" {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
}
func ValidUsernameHandler(w http.ResponseWriter, r *http.Request, hub *websocket.Hub) {
	if r.Method != "POST" {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
}

func ChatHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != "GET" {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	log.Println(vars["room"])

	var chat []database.Message
	db.Where("room-code = ?", vars["room"]).Find(&chat)
	log.Println(chat)

	jsonRes, _ := json.Marshal(chat)
	w.Write([]byte(jsonRes))
}
