package longpoll

import (
	"encoding/json"

	"github.com/Toffee-iZt/HwBot/vkapi"
)

// Event struct.
type Event struct {
	Object interface{}
	Type   string
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (e *Event) UnmarshalJSON(data []byte) error {
	var update struct {
		Type string          `json:"type"`
		Raw  json.RawMessage `json:"object"`
	}
	err := json.Unmarshal(data, &update)
	if err != nil {
		return err
	}

	e.Type = update.Type

	switch e.Type {
	case TypeMessageNew:
		e.Object = new(MessageNew)
	case TypeMessageEvent:
		e.Object = new(MessageEvent)
	default:
		e.Object = update.Raw
		return nil
	}

	return json.Unmarshal(update.Raw, e.Object)
}

//
const (
	TypeMessageNew   = "message_new"
	TypeMessageEvent = "message_event"
)

// MessageNew struct.
type MessageNew struct {
	Message    vkapi.Message `json:"message"`
	ClientInfo struct {
		ButtonActions  []string `json:"button_actions"`
		Keyboard       bool     `json:"keyboard"`
		InlineKeyboard bool     `json:"inline_keyboard"`
		Carousel       bool     `json:"carousel"`
		LangID         int      `json:"lang_id"`
	} `json:"client_info"`
}

// MessageEvent struct.
type MessageEvent struct {
	Payload       interface{}  `json:"payload"`
	EventID       string       `json:"event_id"`
	ConvMessageID int          `json:"conversation_message_id"`
	UserID        vkapi.UserID `json:"user_id"`
	PeerID        vkapi.ID     `json:"peer_id"`
}
