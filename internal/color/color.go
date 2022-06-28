package color

import (
	"fmt"
	"math/rand"
	"time"
)

func Grey(s string) string {
	return "\x1b[38;5;243m" + s + "\x1b[0m"
}

func Custom(s string, color int) string {
	return fmt.Sprintf("\x1b[38;5;%dm", color) + s + "\x1b[0m"
}
func CustomWithBg(s string, color int) string {
	return "\x1b[48;5;239m" + fmt.Sprintf("\x1b[38;5;%dm", color) + s + "\x1b[0m"
}

func Random() int {
	rand.Seed(time.Now().UnixNano())

	return rand.Intn(228) + 1
}
