package longpoll

import (
	"HwBot/vkapi/vktypes"
	"encoding/json"
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
	e.Object = vktypes.Alloc(e.Type)
	if e.Object == nil {
		e.Object = update.Raw
		return nil
	}

	return json.Unmarshal(update.Raw, e.Object)
}

func init() {
	vktypes.Reg(TypeMessageNew, (*MessageNew)(nil))
	vktypes.Reg(TypeMessageEvent, (*MessageEvent)(nil))
}

//
const (
	TypeMessageNew   = "message_new"
	TypeMessageEvent = "message_event"
)

// MessageNew struct.
type MessageNew struct {
	Message    vktypes.Message `json:"message"`
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
	Payload       interface{} `json:"payload"`
	EventID       string      `json:"event_id"`
	ConvMessageID int         `json:"conversation_message_id"`
	UserID        int         `json:"user_id"`
	PeerID        int         `json:"peer_id"`
}
