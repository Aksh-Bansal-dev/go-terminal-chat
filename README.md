# Terminal-Chat
![Screenshot from 2022-06-26 00-07-55](https://user-images.githubusercontent.com/63552235/175786432-8eda8517-0630-4394-ab43-587f499e67b2.png)


A terminal chat app built using Golang.

## How to use?

- Start the server using `go run main.go [-addr]`
- Start the client using `go run client.go [-addr] [-user]`

## How it works?

It uses websockets for server-client communication. The server broadcasts all messages.

## How can I add more emojis?

Add emojis which you want to add in `internal/emoji/emojiMap.go` and create a PR.
