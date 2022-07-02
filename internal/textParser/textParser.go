package textParser

import (
	"regexp"
	"strings"
)

var (
	urlRegex = `https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`
)

func Parse(text string) string {
	arr := strings.Split(text, " ")
	parseEmoji(&arr)
	parseUrl(&arr)
	return strings.Join(arr, " ")
}

func parseUrl(arr *[]string) {
	for i, word := range *arr {
		if match, _ := regexp.MatchString(urlRegex, word); match {
			(*arr)[i] = "🔗\x1b]8;;" + word + "\x07" + word + "\x1b]8;;\x07"
		}
	}
}
