package api

import (
	"HwBot/vkapi/longpoll"
)

type IncomingMessage struct {
	*longpoll.MessageNew
}

type CallbackMessage struct {
	*longpoll.MessageEvent
}
