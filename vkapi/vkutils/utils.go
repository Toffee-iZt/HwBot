package vkutils

import (
	"strconv"
	"strings"

	"github.com/Toffee-iZt/HwBot/vkapi"
)

// Mention returns vk string mention with text.
func Mention(id vkapi.ID, text string) string {
	//if id < 0 {
	//	return "@public" + itoa(-id) + "(" + text + ")"
	//}
	//return "@id" + itoa(id) + "(" + text + ")"

	if id < 0 {
		return "[club" + strconv.Itoa(int(id.ToGroup())) + "|" + text + "]"
	}
	return "[id" + strconv.Itoa(int(id.ToUser())) + "|" + text + "]"
}

// ParseMention extract user id and text from vk mention.
func ParseMention(mention string) (vkapi.ID, string) {
	if mention[0] != '[' || mention[len(mention)-1] != ']' {
		return 0, ""
	}

	di := strings.IndexByte(mention, '|')
	if di == -1 {
		return 0, ""
	}
	id, text := mention[1:di], mention[di:len(mention)-1]

	u := strings.HasPrefix(id, "id")
	if u {
		id = id[2:]
	} else {
		id = id[4:]
	}
	idint, err := strconv.Atoi(id)
	if err != nil {
		return 0, ""
	}

	if !u {
		idint = -idint
	}

	return vkapi.ID(idint), text
}
