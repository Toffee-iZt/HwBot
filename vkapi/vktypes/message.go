package vktypes

// Message ...
type Message struct {
	ID            int          `json:"id"`
	Date          int64        `json:"date"`
	FromID        int          `json:"from_id"`
	PeerID        int          `json:"peer_id"`
	Text          string       `json:"text"`
	Attachments   []Attachment `json:"attachments"`
	Forward       []shortMsg   `json:"fwd_messages"`
	Reply         *shortMsg    `json:"reply_message"`
	ConvMessageID int          `json:"conversation_message_id"`
	Out           int          `json:"out"`
	Important     bool         `json:"important"`
	IsHidden      bool         `json:"is_hidden"`
	Action        *struct {
		Type     string `json:"type"`
		MemberID int    `json:"member_id"`
	} `json:"action"`
}

// Actions types.
const (
	ActionTypeKickUser   = "chat_kick_user"
	ActionTypeInviteUser = "chat_invite_user"
	ActionTypeByLink     = "chat_invite_user_by_link"
)

type shortMsg struct {
	ID            int          `json:"id"`
	Date          int64        `json:"date"`
	FromID        int          `json:"from_id"`
	PeerID        int          `json:"peer_id"`
	Text          string       `json:"text"`
	Attachments   []Attachment `json:"attachments"`
	ConvMessageID int          `json:"conversation_message_id"`
}
