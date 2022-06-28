package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Aksh-Bansal-dev/go-terminal-chat/internal/color"
	"github.com/Aksh-Bansal-dev/go-terminal-chat/internal/emoji"
	"github.com/Aksh-Bansal-dev/go-terminal-chat/internal/tui"
	"github.com/Aksh-Bansal-dev/go-terminal-chat/internal/user"
	"github.com/gorilla/websocket"
)

func printMsg(msg tui.Message) {
	content := emoji.ParseText(msg.Content)
	time := color.Grey(msg.Time)
	if msg.Content == "" {
		return
	} else if err := user.IsValidUsername(msg.Username); err != nil {
		fmt.Printf("%s %s\n", time, content)
	} else if msg.Username == *username {
		fmt.Printf("%s %s: %s\n", time, color.CustomWithBg(msg.Username, msg.Color), content)
	} else {
		fmt.Printf("%s %s: %s\n", time, color.Custom(msg.Username, msg.Color), content)
	}
}

func sendAnnouncement(username string, announcementType string, writeMessage func(messageType int, data []byte) error) error {
	s := username + " " + announcementType
	newMsg := tui.NewMessage(s+" the chat!", s, myColor)
	postBody, _ := json.Marshal(newMsg)
	err := writeMessage(websocket.TextMessage, []byte(postBody))
	return err
}

func sendMsg(content string, writeMessage func(messageType int, data []byte) error) error {
	newMsg := tui.NewMessage(content, *username, myColor)
	postBody, _ := json.Marshal(newMsg)
	err := writeMessage(websocket.TextMessage, []byte(postBody))
	return err
}

func sendMessageToServer(done <-chan struct{}, interrupt <-chan os.Signal, input <-chan string, c websocket.Conn) {
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			fmt.Println("interrupt")

			err := sendAnnouncement(*username, "left", c.WriteMessage)
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

func readSTDIN(input chan string, reader *bufio.Reader) {
	defer close(input)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		input <- text
	}
}

func readMessageFromServer(done chan struct{}, c websocket.Conn) {
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
}
