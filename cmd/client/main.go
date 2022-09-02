package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"regexp"

	"github.com/Aksh-Bansal-dev/go-terminal-chat/internal/color"
	"github.com/Aksh-Bansal-dev/go-terminal-chat/internal/tui"
	"github.com/Aksh-Bansal-dev/go-terminal-chat/internal/user"
	"github.com/gorilla/websocket"
)

var (
	addr = flag.String("addr", "localhost:8080", `http service address
EXAMPLE: ./client -addr localhost:5000	
	`)
	username = flag.String("user", "Newbie", `username for chat
EXAMPLE: ./client -user aksh	
	`)
	tuiMode = flag.Bool("tui", false, `Use this app in tui mode
EXAMPLE: ./client -tui	
	`)
	roomCode = flag.String("room", "general", `Specify room code
EXAMPLE: ./client -room private	
	`)
	myColor int
)

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile)
	if err := user.IsValidUsername(*username, *addr, *roomCode); err != nil {
		fmt.Println(err)
		return
	}
	roomCodeR, _ := regexp.Compile("[a-z]{3,6}")
	if !roomCodeR.MatchString(*roomCode) {
		fmt.Println("Room code must be 3 to 6 letter string with lowercase alphabets(a-z)")
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

	user.GetInitialUsers(*addr, *roomCode)
	done := make(chan struct{})
	input := make(chan string)

	if *tuiMode {
		go tui.Run(*username, &input, *addr)
		for _, msg := range user.GetChat(*addr, *roomCode) {
			tui.PrintMessage(msg)
		}
	} else {
		for _, msg := range user.GetChat(*addr, *roomCode) {
			printMsg(msg)
		}
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
