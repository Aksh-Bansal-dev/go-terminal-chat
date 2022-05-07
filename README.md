# Terminal-Chat
A simple terminal chat app built using Golang.

> Note: This app uses no external golang package

## How to use?
- Start the server using `go run main.go [-addr]`
- Start the client using `go run client.go [-addr] [-user]`

## How it works?

The clients check server regularly and pull new messages if they are present. 
The server stores all the messages in memory(no database).
