package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/Aksh-Bansal-dev/go-terminal-chat/internal/color"
	"github.com/Aksh-Bansal-dev/go-terminal-chat/internal/tui"
	"github.com/Aksh-Bansal-dev/go-terminal-chat/internal/user"
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
	log.SetFlags(log.Lshortfile)
	if err := user.IsValidUsername(*username, *addr); err != nil {
		fmt.Println(err)
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

	user.GetInitialUsers(*addr)
	done := make(chan struct{})
	input := make(chan string)

	if *tuiMode {
		go tui.Run(*username, &input)
	}

	// Welcome message
	err = sendAnnouncement(*username, "joined", c.WriteMessage)
	if err != nil {
		log.Println("err:", err)
		return
	}

	// Read messages
	go readMessageFromServer(done, *c)

	if !*tuiMode {
		// Take input from STDIN
		reader := bufio.NewReader(os.Stdin)
		go readSTDIN(input, reader)
	}

	// Send to message server
	defer close(input)
	sendMessageToServer(done, interrupt, input, *c)
}
