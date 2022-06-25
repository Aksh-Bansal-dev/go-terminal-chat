package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/Aksh-Bansal-dev/go-terminal-chat/internal/color"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var username = flag.String("user", "Newbie", "username for chat")

type Message struct {
	Username string
	Content  string
	Time     string
	Color    string
}

var myColor string

func main() {
	flag.Parse()

	myColor = color.Random()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	reader := bufio.NewReader(os.Stdin)
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	done := make(chan struct{})
	input := make(chan string)

	// Read messages
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			var msg Message
			err = json.Unmarshal(message, &msg)
			if err != nil {
				fmt.Println("err")
				return
			}
			fmt.Printf("%s %s: %s\n", color.Grey(msg.Time), color.Custom(msg.Username, msg.Color), msg.Content)
		}
	}()

	// Take input from STDIN
	go func() {
		defer close(input)
		for {
			text, err := reader.ReadString('\n')
			if err != nil {
				continue
			}
			input <- text
		}
	}()

	// Send to message server
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			fmt.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Println("write close:", err)
				return
			}
			os.Exit(0)
		case text := <-input:
			t := time.Now()
			newMsg := Message{
				Username: *username,
				Content:  text[:len(text)-1],
				Time:     fmt.Sprintf("%d:%d:%d", t.Hour(), t.Minute(), t.Second()),
				Color:    myColor,
			}
			postBody, _ := json.Marshal(newMsg)
			err := c.WriteMessage(websocket.TextMessage, []byte(postBody))
			if err != nil {
				fmt.Println("write:", err)
				return
			}
		}
	}
}
