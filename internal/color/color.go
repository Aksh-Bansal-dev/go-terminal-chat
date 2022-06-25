package color

import (
	"fmt"
	"math/rand"
	"time"
)

func Grey(s string) string {
	return "\x1b[38;5;243m" + s + "\x1b[0m"
}

func Custom(s string, color string) string {
	return color + s + "\x1b[0m"
}

func Random() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("\x1b[38;5;%dm", rand.Intn(228)+1)
}
