package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"

	"github.com/Aksh-Bansal-dev/go-terminal-chat/internal/color"
	"github.com/Aksh-Bansal-dev/go-terminal-chat/internal/emoji"
	"github.com/Aksh-Bansal-dev/go-terminal-chat/internal/tui"
	"github.com/gorilla/websocket"
)

var (
	addr     = flag.String("addr", "localhost:8080", "http service address")
	username = flag.String("user", "Newbie", "username for chat")
	tuiMode  = flag.Bool("tui", false, "Use this app in tui mode")
	myColor  int
)

func main() {
	flag.Parse()
	if *username == "" {
		fmt.Println("Username cannot be empty")
		return
	}

	myColor = color.Random()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	done := make(chan struct{})
	input := make(chan string)

	if *tuiMode {
		go tui.Run(*username, &input)
	}

	// Welcome message
	err = sendAnnouncement(*username+" joined the chat!", c.WriteMessage)
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	// Read messages
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				fmt.Println("websocket error: ", err)
				return
			}
			var msg tui.Message
			err = json.Unmarshal(message, &msg)
			if err != nil {
				fmt.Println("Parsing error:", err)
				return
			}
			if *tuiMode {
				tui.PrintMessage(msg)
			} else {
				printMsg(msg)
			}
		}
	}()

	if !*tuiMode {
		reader := bufio.NewReader(os.Stdin)
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
	}

	// Send to message server
	defer close(input)
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
			if text[len(text)-1] == '\n' {
				text = text[:len(text)-1]
			}
			err := sendMsg(text, c.WriteMessage)
			if err != nil {
				fmt.Println("write:", err)
				return
			}
		}
	}
}

func sendMsg(content string, writeMessage func(messageType int, data []byte) error) error {
	newMsg := tui.NewMessage(content, *username, myColor)
	postBody, _ := json.Marshal(newMsg)
	err := writeMessage(websocket.TextMessage, []byte(postBody))
	return err
}

func sendAnnouncement(content string, writeMessage func(messageType int, data []byte) error) error {
	newMsg := tui.NewMessage(content, "", myColor)
	postBody, _ := json.Marshal(newMsg)
	err := writeMessage(websocket.TextMessage, []byte(postBody))
	return err
}

func printMsg(msg tui.Message) {
	content := emoji.ParseText(msg.Content)
	time := color.Grey(msg.Time)
	if msg.Content == "" {
		return
	} else if msg.Username == "" {
		fmt.Printf("%s %s\n", time, content)
	} else if msg.Username == *username {
		fmt.Printf("%s %s: %s\n", time, color.CustomWithBg(msg.Username, msg.Color), content)
	} else {
		fmt.Printf("%s %s: %s\n", time, color.Custom(msg.Username, msg.Color), content)
	}
}
