package vkapi

import (
	"encoding/json"
	"sync/atomic"
)

// OutMessage struct.
type OutMessage struct {
	UserID  UserID
	UserIDs []UserID
	PeerID  ID
	PeerIDs []ID
	Domain  string
	ChatID  ChatID

	OutMessageContent
}

// OutMessageContent struct.
type OutMessageContent struct {
	Message    string
	Lat, Long  float64
	Attachment []string
	StickerID  int
	Keyboard   *Keyboard
	//Template map[string]interface{}
	//Payload       map[string]interface{} // ???
	//ContentSource map[string]interface{} // ???

	ReplyTo         int
	ForwardMessages []int
	//Forward common.JSONData

	DontParseLinks  bool
	DisableMentions bool
}

// Send sends a message.
func (c *Client) Send(msg OutMessage) (int, error) {
	if msg.Message == "" && msg.Attachment == nil {
		return 0, nil
	}

	if len(msg.Attachment) > 10 {
		msg.Attachment = msg.Attachment[:10]
	}

	args := vkargs{
		"random_id": atomic.AddInt32(&c.rndID, 1),
		"user_id":   msg.UserID,
		"user_ids":  msg.UserIDs,
		"peer_id":   msg.PeerID,
		"peer_ids":  msg.PeerIDs,
		"domain":    msg.Domain,
		"chat_id":   msg.ChatID,

		"message":    msg.Message,
		"lat":        msg.Lat,
		"long":       msg.Long,
		"attachment": msg.Attachment,
		"sticker_id": msg.StickerID,
		"keyboard":   msg.Keyboard.Data(),

		"reply_to":         msg.ReplyTo,
		"forward_messages": msg.ForwardMessages,
	}

	if msg.DontParseLinks {
		args["dont_parse_links"] = "1"
	}
	if msg.DisableMentions {
		args["disable_mentions"] = "1"
	}

	var id int
	return id, c.method(&id, "messages.send", args)
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
func (c *Client) SendMessageEventAnswer(eventID string, userID UserID, peerID ID, eventData *EventData) error {
	var data string
	if eventData != nil {
		m, _ := json.Marshal(eventData)
		data = string(m)
	}

	return c.method(nil, "messages.sendMessageEventAnswer", vkargs{
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
func (c *Client) GetConversationMembers(peerID ID) (*Members, error) {
	var mem Members
	return &mem, c.method(&mem, "messages.getConversationMembers", vkargs{
		"peer_id": peerID,
	})
}

// RemoveChatUser allows the current user to leave a chat or, if the current user started the chat,
// allows the user to remove another user from the chat.
func (c *Client) RemoveChatUser(chatID ChatID, memberID ID) error {
	return c.method(nil, "messages.removeChatUser", vkargs{
		"chat_id":   chatID,
		"member_id": memberID,
	})
}
