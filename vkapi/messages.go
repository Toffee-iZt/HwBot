package vkapi

import (
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

	args := newArgs()
	args.Set("random_id", itoa(int(atomic.AddInt32(&c.rndID, 1))))

	args.Set("user_id", itoa(int(msg.UserID)))
	for i := range msg.UserIDs {
		args.Add("user_ids", itoa(int(msg.UserIDs[i])))
	}
	args.Set("peer_id", itoa(int(msg.PeerID)))
	for i := range msg.PeerIDs {
		args.Add("peer_ids", itoa(int(msg.PeerIDs[i])))
	}
	args.Set("domain", msg.Domain)
	args.Set("chat_id", itoa(int(msg.ChatID)))

	args.Set("message", msg.Message)
	args.Set("lat", ftoa(msg.Lat))
	args.Set("long", ftoa(msg.Long))
	args.Set("attachment", msg.Attachment...)
	args.Set("sticker_id", itoa(msg.StickerID))
	args.Set("keyboard", msg.Keyboard.String())

	args.Set("reply_to", itoa(msg.ReplyTo))
	for _, f := range msg.ForwardMessages {
		args.Add("peer_ids", itoa(f))
	}

	if msg.DontParseLinks {
		args.Set("dont_parse_links", "1")
	}
	if msg.DisableMentions {
		args.Set("disable_mentions", "1")
	}

	var id int
	return id, c.method("messages.send", args, &id)
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
	args := newArgs()
	args.Set("event_id", eventID)
	args.Set("user_id", itoa(int(userID)))
	args.Set("peer_id", itoa(int(peerID)))
	if eventData != nil {
		args.Set("event_data", string(marshal(eventData)))
	}

	return c.method("messages.sendMessageEventAnswer", args, nil)
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
	args := newArgs()
	args.Set("peer_id", itoa(int(peerID)))
	var mem Members
	return &mem, c.method("messages.getConversationMembers", args, &mem)
}

// RemoveChatUser allows the current user to leave a chat or, if the current user started the chat,
// allows the user to remove another user from the chat.
func (c *Client) RemoveChatUser(chatID ChatID, memberID ID) error {
	args := newArgs()
	args.Set("chat_id", itoa(int(chatID)))
	args.Set("member_id", itoa(int(memberID)))
	return c.method("messages.removeChatUser", args, nil)
}
