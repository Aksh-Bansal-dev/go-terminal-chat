package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var username = flag.String("user", "Newbie", "username for chat")

type Message struct {
	Name    string
	Content string
}

var idx int
var m sync.Mutex

func main() {
	flag.Parse()
	idx = 0
	newMsg := Message{
		Name:    *username,
		Content: fmt.Sprintf("%s joined the chat!", *username),
	}
	postBody, _ := json.Marshal(newMsg)
	initialBody := bytes.NewBuffer(postBody)
	_, err := http.Post(fmt.Sprintf("http://%s/sendMsg", *addr), "application/json", initialBody)
	if err != nil {
		panic(err)
	}
	go sendMsg()
	go recvMsg()
	select {}
}

func recvMsg() {
	for {
		m.Lock()
		res, err := http.Get(fmt.Sprintf("http://%s/getMsg?idx=%d", *addr, idx))
		if err != nil {
			panic(err)
		}
		var allMsg []Message
		err2 := json.NewDecoder(res.Body).Decode(&allMsg)
		if err2 != nil {
			return
		}
		idx += len(allMsg)
		for i := 0; i < len(allMsg); i++ {
			fmt.Println(allMsg[i].Content)
		}
		m.Unlock()
		time.Sleep(time.Second)
	}
}

func sendMsg() {
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')

		newMsg := Message{
			Name:    *username,
			Content: fmt.Sprintf("[%s] %s", *username, text[:len(text)-1]),
		}
		postBody, _ := json.Marshal(newMsg)
		body := bytes.NewBuffer(postBody)
		m.Lock()
		_, err := http.Post(fmt.Sprintf("http://%s/sendMsg", *addr), "application/json", body)
		idx++
		m.Unlock()
		if err != nil {
			panic(err)
		}
	}
}
