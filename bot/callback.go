package bot

import (
	"encoding/json"

	"github.com/Toffee-iZt/HwBot/vkapi"
	"github.com/Toffee-iZt/HwBot/vkapi/longpoll"
)

func (b *Bot) onCallback(callback *longpoll.MessageEvent) {
	modname, data := unwrapPayload(callback.Payload)
	if modname == "" {
		return
	}
	var mod *Module
	for _, mod = range b.modules {
		if modname == mod.Name {
			break
		}
	}
	if mod.Callback == nil {
		return
	}
	ectx := &CallbackContext{
		eventctx: eventctx{
			Conv: &Conversation{
				peer: callback.PeerID,
				api:  b.vk,
			},
			mod: mod,
		},
		eventID: callback.EventID,
		userID:  callback.UserID,
		msgID:   callback.ConversationMessageID,
	}
	mod.Callback(ectx, data)
}

func wrapPayload(mod *Module, payload interface{}) vkapi.JSONData {
	var botPayload = struct {
		HwBot string      `json:"hwbot"`
		Data  interface{} `json:"data"`
	}{
		HwBot: mod.Name,
		Data:  payload,
	}

	d, _ := vkapi.NewJSONDataMarshal(botPayload)
	return d
}

func unwrapPayload(payload vkapi.JSONData) (string, vkapi.JSONData) {
	var botPayload struct {
		HwBot string         `json:"hwbot"`
		Data  vkapi.JSONData `json:"data"`
	}
	err := json.Unmarshal([]byte(payload), &botPayload)
	if err != nil {
		return "", ""
	}
	return botPayload.HwBot, botPayload.Data
}
