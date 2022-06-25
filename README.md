# Terminal-Chat

A simple multithreaded terminal chat app built using Golang.

## How to use?

- Start the server using `go run main.go [-addr]`
- Start the client using `go run client.go [-addr] [-user]`

## How it works?

It uses websockets for server-client communication. The server broadcasts all messages.

## How can I add more emojis?

Add emojis which you want to add in `internal/emoji/emojiMap.go` and create a PR.
