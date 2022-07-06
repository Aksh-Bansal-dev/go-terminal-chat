# Terminal-Chat

An awesome terminal chat application built using Golang.

### TUI mode

![Screenshot from 2022-06-30 22-34-22](https://user-images.githubusercontent.com/63552235/176736275-298b4876-5bec-4ff6-9f9f-55270be0cdd7.png)

### CLI mode

![Screenshot from 2022-06-30 22-34-58](https://user-images.githubusercontent.com/63552235/176736282-9c9b18db-dd8a-4423-8b2e-62b53822972a.png)

> Web client coming soon.

## How to use?

- Start the server using `go run cmd/server/main.go [-addr]`
- Start the client using `go run cmd/client/main.go [-addr] [-user] [-tui]`

Don't have golang installed? Download the executables from [here](https://github.com/Aksh-Bansal-dev/go-terminal-chat/releases/tag/v1.0.0).

- Use `:<emoji-code>:` to send emoji.
- Use `>username <msg>` to send a private message to username.

## How it works?

It uses websockets for server-client communication. The server broadcasts all messages.

## How can I add more emojis?

Add emojis which you want to add in `internal/textParser/emoji.go` and create a PR.
