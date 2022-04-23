package vkapi

import (
	"encoding/json"
	"sync/atomic"
)

// OutMessage struct.
type OutMessage struct {
	randomID int32 `vkargs:"random_id"`

	UserID  UserID   `vkargs:"user_id"`
	UserIDs []UserID `vkargs:"user_ids"`
	PeerID  ID       `vkargs:"peer_id"`
	PeerIDs []ID     `vkargs:"peer_ids"`
	Domain  string   `vkargs:"domain"`
	ChatID  ChatID   `vkargs:"chat_id"`

	OutMessageContent
}

// OutMessageContent struct.
type OutMessageContent struct {
	Message    string   `vkargs:"message"`
	Lat        float64  `vkargs:"lat"`
	Long       float64  `vkargs:"long"`
	Attachment []string `vkargs:"attachment"`
	StickerID  int      `vkargs:"sticker_id"`
	Keyboard   JSONData `vkargs:"keyboard"`

	ReplyTo         int   `vkargs:"reply_to"`
	ForwardMessages []int `vkargs:"forward_messages"`

	DontParseLinks  bool `vkargs:"dont_parse_links"`
	DisableMentions bool `vkargs:"disable_mentions"`
}

// SendMessage sends a message.
func (c *Client) SendMessage(msg OutMessage) *Error {
	msg.randomID = atomic.AddInt32(&c.rndID, 1)
	return c.method(nil, "messages.send", msg)
}

// EditMessage edits a message.
func (c *Client) EditMessage(peerID ID, convMessageID int, content OutMessageContent) (bool, *Error) {
	var edit = struct {
		PeerID                ID  `vkargs:"peer_id"`
		ConversationMessageID int `vkargs:"conversation_message_id"`

		OutMessageContent
	}{
		PeerID:                peerID,
		ConversationMessageID: convMessageID,
		OutMessageContent:     content,
	}

	var ok boolean
	return bool(ok), c.method(&ok, "messages.edit", edit)
}

//
const (
	EventDataTypeShowSnackbar = "show_snackbar"
	EventDataTypeOpenLink     = "open_link"
	EventDataTypeOpenApp      = "open_app"
)

// EventData struct.
type EventData struct {
	Type string `json:"type"`

	// show_snackbar
	Text string `json:"text,omitempty"`

	// open_link
	Link string `json:"link,omitempty"`

	// open_app
	Hash    string `json:"hash,omitempty"`
	AppID   int    `json:"app_id,omitempty"`
	OwnerID int    `json:"owner_id,omitempty"`
}

// SendMessageEventAnswer sends event answer.
func (c *Client) SendMessageEventAnswer(eventID string, userID UserID, peerID ID, eventData *EventData) *Error {
	var data string
	if eventData != nil {
		m, _ := json.Marshal(eventData)
		data = string(m)
	}

	return c.method(nil, "messages.sendMessageEventAnswer", ArgsMap{
		"event_id":   eventID,
		"user_id":    userID,
		"peer_id":    peerID,
		"event_data": data,
	})
}

// Members struct.
type Members struct {
	Items    []*Member `json:"items"`
	Profiles []*User   `json:"profiles"`
	Groups   []*Group  `json:"groups"`
	Count    int       `json:"count"`
}

// Member struct.
type Member struct {
	MemberID  ID    `json:"member_id"`
	InvitedBy ID    `json:"invited_by"`
	JoinDate  int64 `json:"join_date"`
	IsAdmin   bool  `json:"is_admin"`
	IsOwner   bool  `json:"is_owner"`
	CanKick   bool  `json:"can_kick"`
}

// GetConversationMembers returns a list of IDs of users participating in a conversation.
func (c *Client) GetConversationMembers(peerID ID) (*Members, *Error) {
	var mem Members
	return &mem, c.method(&mem, "messages.getConversationMembers", ArgsMap{
		"peer_id": peerID,
	})
}

// RemoveChatUser allows the current user to leave a chat or, if the current user started the chat,
// allows the user to remove another user from the chat.
func (c *Client) RemoveChatUser(chatID ChatID, memberID ID) *Error {
	return c.method(nil, "messages.removeChatUser", ArgsMap{
		"chat_id":   chatID,
		"member_id": memberID,
	})
}

// Message struct.
type Message struct {
	Date          int64        `json:"date"`
	FromID        ID           `json:"from_id"`
	PeerID        ID           `json:"peer_id"`
	Text          string       `json:"text"`
	Attachments   []attachment `json:"attachments"`
	Forward       []shortMsg   `json:"fwd_messages"`
	Reply         *shortMsg    `json:"reply_message"`
	Payload       JSONData     `json:"payload"`
	ConvMessageID int          `json:"conversation_message_id"`
	Action        *struct {
		Type     string `json:"type"`
		MemberID int    `json:"member_id"`
	} `json:"action"`
}

type shortMsg struct {
	ID            int          `json:"id"`
	Date          int64        `json:"date"`
	FromID        ID           `json:"from_id"`
	PeerID        ID           `json:"peer_id"`
	Text          string       `json:"text"`
	Attachments   []attachment `json:"attachments"`
	ConvMessageID int          `json:"conversation_message_id"`
}

type attachment struct {
	Photo        *Photo        `json:"photo"`
	AudioMessage *AudioMessage `json:"audio_message"`
	Type         string        `json:"type"`
}
