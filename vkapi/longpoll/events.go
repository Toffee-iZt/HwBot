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
	Message    Message    `json:"message"`
	ClientInfo ClientInfo `json:"client_info"`
}

// MessageEvent struct.
type MessageEvent struct {
	Payload       vkapi.JSONData `json:"payload"`
	EventID       string         `json:"event_id"`
	ConvMessageID int            `json:"conversation_message_id"`
	UserID        vkapi.UserID   `json:"user_id"`
	PeerID        vkapi.ID       `json:"peer_id"`
}

// Message struct.
type Message struct {
	ID            int                `json:"id"`
	Date          int64              `json:"date"`
	FromID        vkapi.ID           `json:"from_id"`
	PeerID        vkapi.ID           `json:"peer_id"`
	Text          string             `json:"text"`
	Attachments   []vkapi.Attachment `json:"attachments"`
	Forward       []shortMsg         `json:"fwd_messages"`
	Reply         *shortMsg          `json:"reply_message"`
	ConvMessageID int                `json:"conversation_message_id"`
	Out           int                `json:"out"`
	Important     bool               `json:"important"`
	IsHidden      bool               `json:"is_hidden"`
	Action        *struct {
		Type     string `json:"type"`
		MemberID int    `json:"member_id"`
	} `json:"action"`
}

// ClientInfo struct.
type ClientInfo struct {
	ButtonActions  []string `json:"button_actions"`
	Keyboard       bool     `json:"keyboard"`
	InlineKeyboard bool     `json:"inline_keyboard"`
	Carousel       bool     `json:"carousel"`
	LangID         int      `json:"lang_id"`
}

// Actions types.
const (
	ActionTypeKickUser   = "chat_kick_user"
	ActionTypeInviteUser = "chat_invite_user"
	ActionTypeByLink     = "chat_invite_user_by_link"
)

type shortMsg struct {
	ID            int                `json:"id"`
	Date          int64              `json:"date"`
	FromID        vkapi.ID           `json:"from_id"`
	PeerID        vkapi.ID           `json:"peer_id"`
	Text          string             `json:"text"`
	Attachments   []vkapi.Attachment `json:"attachments"`
	ConvMessageID int                `json:"conversation_message_id"`
}
