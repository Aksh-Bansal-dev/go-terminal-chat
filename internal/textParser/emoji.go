package textParser

var emojiMap = map[string]string{
	":joy:":        "😂",
	":+1:":         "👍",
	":thumbsup:":   "👍",
	":thumbsdown:": "👎",
	":-1:":         "👎",
	":fire:":       "🔥",
	":emotional:":  "🥺",
	":cry:":        "😢",
	":poop:":       "💩",
}

func parseEmoji(arr *[]string) {
	for i, word := range *arr {
		if val, ok := emojiMap[word]; ok {
			(*arr)[i] = val
		}
	}
}
