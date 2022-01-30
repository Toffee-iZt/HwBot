package vkapi

import (
	"encoding/json"
	"sync/atomic"
)

// OutMessage struct.
type OutMessage struct {
	randomID int32 `vkargs:"random_id"`

	UserID  UserID   `vkargs:"user_id,omitempty"`
	UserIDs []UserID `vkargs:"user_ids,omitempty"`
	PeerID  ID       `vkargs:"peer_id,omitempty"`
	PeerIDs []ID     `vkargs:"peer_ids,omitempty"`
	Domain  string   `vkargs:"domain,omitempty"`
	ChatID  ChatID   `vkargs:"chat_id,omitempty"`

	OutMessageContent `vkargs:"embed"`
}

// OutMessageContent struct.
type OutMessageContent struct {
	Message    string   `vkargs:"message,omitempty"`
	Lat        float64  `vkargs:"lat,omitempty"`
	Long       float64  `vkargs:"long,omitempty"`
	Attachment []string `vkargs:"attachment,omitempty"`
	StickerID  int      `vkargs:"sticker_id,omitempty"`
	Keyboard   JSONData `vkargs:"keyboard,omitempty"`

	ReplyTo         int   `vkargs:"reply_to,omitempty"`
	ForwardMessages []int `vkargs:"forward_messages,omitempty"`

	DontParseLinks  bool `vkargs:"dont_parse_links,omitempty"`
	DisableMentions bool `vkargs:"disable_mentions,omitempty"`
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

		OutMessageContent `vkargs:"embed"`
	}{
		PeerID:                peerID,
		ConversationMessageID: convMessageID,
		OutMessageContent:     content,
	}

	var ok BoolInt
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
