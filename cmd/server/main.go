package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

type Message struct {
	Name    string
	Content string
}

var allMsg []Message

func getMsg(w http.ResponseWriter, r *http.Request) {
	idx, _ := strconv.Atoi(r.URL.Query().Get("idx"))
	w.Header().Set("Content-Type", "application/json")
	if idx >= len(allMsg) {
		json.NewEncoder(w).Encode(allMsg[0:0])
	} else {
		json.NewEncoder(w).Encode(allMsg[idx:])
	}
}

func sendMsg(w http.ResponseWriter, r *http.Request) {
	var msg Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	allMsg = append(allMsg, msg)
}

func main() {
	flag.Parse()
	http.HandleFunc("/getMsg", getMsg)
	http.HandleFunc("/sendMsg", sendMsg)
	log.Println(fmt.Sprintf("Server started at %s", *addr))
	log.Fatal(http.ListenAndServe(*addr, nil))
}
