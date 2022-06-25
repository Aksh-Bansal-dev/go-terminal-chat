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

func newMessage(content string, customUserName string) Message {
	t := time.Now()
	return Message{
		Username: customUserName,
		Content:  content,
		Color:    myColor,
		Time:     fmt.Sprintf("%d:%d:%d", t.Hour(), t.Minute(), t.Second()),
	}
}

var myColor string

func main() {
	flag.Parse()
	if *username == "" {
		fmt.Println("Username cannot be empty")
		return
	}

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

	// Welcome message
	err = sendAnnouncement(*username+" joined the chat!", c.WriteMessage)
	if err != nil {
		fmt.Println("err:", err)
		return
	}

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
			printMsg(msg)
		}
	}()

	// Take input from STDIN
	go func() {
		defer close(input)
		for {
			// printMsg(newMessage("", *username))
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

			err := sendAnnouncement(*username+" left the chat!", c.WriteMessage)
			if err != nil {
				fmt.Println("err:", err)
				return
			}
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Println("write close:", err)
				return
			}
			os.Exit(0)
		case text := <-input:
			err := sendMsg(text[:len(text)-1], c.WriteMessage)
			if err != nil {
				fmt.Println("write:", err)
				return
			}
		}
	}
}

func sendMsg(content string, writeMessage func(messageType int, data []byte) error) error {
	newMsg := newMessage(content, *username)
	postBody, _ := json.Marshal(newMsg)
	err := writeMessage(websocket.TextMessage, []byte(postBody))
	return err
}

func sendAnnouncement(content string, writeMessage func(messageType int, data []byte) error) error {
	newMsg := newMessage(content, "")

	postBody, _ := json.Marshal(newMsg)
	err := writeMessage(websocket.TextMessage, []byte(postBody))
	return err
}

func printMsg(msg Message) {
	if msg.Username == "" {
		fmt.Printf("%s %s\n", color.Grey(msg.Time), msg.Content)
	} else if msg.Username == *username {
		fmt.Printf("%s %s: %s\n", color.Grey(msg.Time), color.CustomWithBg(msg.Username, msg.Color), msg.Content)
	} else if msg.Username == *username && msg.Content == "" {
		fmt.Printf("%s %s: %s ", color.Grey(msg.Time), color.CustomWithBg(msg.Username, msg.Color), msg.Content)
	} else {
		fmt.Printf("%s %s: %s\n", color.Grey(msg.Time), color.Custom(msg.Username, msg.Color), msg.Content)
	}
}
