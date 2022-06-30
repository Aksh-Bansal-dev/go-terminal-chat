package tui

import (
	"fmt"
	"log"
	"syscall"
	"time"

	"github.com/Aksh-Bansal-dev/go-terminal-chat/internal/emoji"
	"github.com/Aksh-Bansal-dev/go-terminal-chat/internal/user"
	"github.com/marcusolsson/tui-go"
)

type Message struct {
	Username string
	Content  string
	Time     string
	Color    int
}

var (
	messages   = []Message{}
	username   string
	input      chan string
	messageBox *tui.Box
	UI         tui.UI
	sidebar    *tui.Box
)

func Run(usrname string, inp *chan string) {
	username = usrname
	input = *inp

	sidebar = tui.NewVBox(
		tui.NewLabel("  	Online  	"),
		tui.NewSpacer(),
	)
	sidebar.SetBorder(true)
	go UpdateSidebar()

	messageBox = tui.NewVBox()

	messageBoxScroll := tui.NewScrollArea(messageBox)
	messageBoxScroll.SetAutoscrollToBottom(true)

	messageBoxContainer := tui.NewVBox(messageBoxScroll)
	messageBoxContainer.SetBorder(true)

	inputBox := tui.NewEntry()
	inputBox.SetFocused(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBoxContainer := tui.NewHBox(inputBox)
	inputBoxContainer.SetBorder(true)
	inputBoxContainer.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(messageBoxContainer, inputBoxContainer)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	inputBox.OnSubmit(func(e *tui.Entry) {
		if e.Text() == "" {
			return
		}
		input <- e.Text()
		inputBox.SetText("")
	})

	root := tui.NewHBox(sidebar, chat)

	ui, err := tui.New(root)
	UI = ui
	if err != nil {
		log.Fatal(err)
	}

	theme := tui.NewTheme()
	theme.SetStyle("label.time", tui.Style{Fg: 243})
	for i := 1; i < 255; i++ {
		theme.SetStyle(fmt.Sprintf("label.color%d", i), tui.Style{Fg: tui.Color(i)})
		theme.SetStyle(fmt.Sprintf("label.wb-color%d", i), tui.Style{Fg: tui.Color(i), Bg: 239})
	}
	ui.SetTheme(theme)

	ui.SetKeybinding("Esc", func() {
		ui.Quit()
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	})
	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}

func PrintMessage(msg Message) {
	tim := tui.NewLabel(msg.Time)
	tim.SetStyleName("time")

	name := tui.NewLabel(msg.Username)
	if err := user.IsValidUsername(msg.Username); err != nil {
		name.SetText("")
	} else if msg.Username == username {
		name.SetStyleName(fmt.Sprintf("wb-color%d", msg.Color))
	} else {
		name.SetStyleName(fmt.Sprintf("color%d", msg.Color))
	}

	content := tui.NewLabel(emoji.ParseText(msg.Content))

	messageBox.Append(tui.NewHBox(
		tim,
		tui.NewPadder(1, 0, name),
		tui.NewPadder(1, 0, content),
		tui.NewSpacer(),
	))
	if UI != nil {
		UI.Update(func() {})
	}
}

func NewMessage(content string, customUserName string, color int) Message {
	t := time.Now()
	return Message{
		Username: customUserName,
		Content:  content,
		Color:    color,
		Time:     fmt.Sprintf("%d:%d:%d", t.Hour(), t.Minute(), t.Second()),
	}
}

func UpdateSidebar() {
	defer close(user.Users)
	userMap := map[string]int{}
	for {
		select {
		case user := <-user.Users:
			for range userMap {
				sidebar.Remove(1)
			}
			if user.Username[0] == ' ' {
				delete(userMap, user.Username[1:])
			} else {
				userMap[user.Username] = user.Color
			}
			for key, color := range userMap {
				userLabel := tui.NewLabel(key)
				if key == username {
					userLabel.SetStyleName(fmt.Sprintf("wb-color%d", color))
				} else {
					userLabel.SetStyleName(fmt.Sprintf("color%d", color))
				}
				sidebar.Insert(1, userLabel)
			}
			if UI != nil {
				UI.Update(func() {})
			}
		}
	}
}
