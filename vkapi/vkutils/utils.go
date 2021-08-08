package vkutils

import (
	"strconv"
	"strings"
)

// Mention returns vk string mention with text.
func Mention(id int, text string) string {
	//if id < 0 {
	//	return "@public" + itoa(-id) + "(" + text + ")"
	//}
	//return "@id" + itoa(id) + "(" + text + ")"

	if id < 0 {
		return "[club" + strconv.Itoa(-id) + "|" + text + "]"
	}
	return "[id" + strconv.Itoa(id) + "|" + text + "]"
}

// ParseMention extract user id and text from vk mention.
func ParseMention(mention string) (int, string) {
	if mention[0] != '[' || mention[len(mention)-1] != ']' {
		return 0, ""
	}
	s := strings.SplitN(mention[1:len(mention)-1], "|", 2)
	if len(s) != 2 {
		return 0, ""
	}

	user := strings.Index(s[0], "id") != -1

	var strID string
	if user {
		strID = strings.Replace(s[0], "id", "", 1)
	} else {
		strID = strings.Replace(s[0], "public", "", 1)
	}

	id, err := strconv.Atoi(strID)
	if err != nil {
		return 0, ""
	}

	if !user {
		id = -id
	}

	return id, s[1]
}
