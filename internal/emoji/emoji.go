package emoji

import "strings"

func ParseText(text string) string {
	arr := strings.Split(text, " ")
	for i, word := range arr {
		if val, ok := emojiMap[word]; ok {
			arr[i] = val
		}
	}
	return strings.Join(arr, " ")
}
