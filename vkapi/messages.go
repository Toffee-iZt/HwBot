package vkapi

import (
	"sync/atomic"
)

// ProvideMessages makes messages provider.
func ProvideMessages(c *Client) *MessagesProvider {
	return &MessagesProvider{
		APIProvider: APIProvider{c},
		randomID:    -(1 << 31),
	}
}

// MessagesProvider provides messages api.
type MessagesProvider struct {
	APIProvider

	randomID int32
}

// OutMessage struct.
type OutMessage struct {
	UserID  UserID
	UserIDs []UserID
	PeerID  ID
	PeerIDs []ID
	Domain  string
	ChatID  ChatID

	Message    string
	Lat, Long  int
	Attachment []string
	StickerID  int
	Keyboard   *Keyboard
	//Template map[string]interface{}
	//Payload       map[string]interface{} // ???
	//ContentSource map[string]interface{} // ???

	ReplyTo         int
	ForwardMessages []int
	//Forward common.JSONData

	DontParseLinks  boolInt
	DisableMentions boolInt
}

// Send sends a message.
func (m *MessagesProvider) Send(msg OutMessage) (int, error) {
	if msg.Message == "" && msg.Attachment == nil {
		return 0, nil
	}

	if len(msg.Attachment) > 10 {
		msg.Attachment = msg.Attachment[:10]
	}

	args := NewArgs()
	args.Set("random_id", itoa(int(atomic.AddInt32(&m.randomID, 1))))

	args.Set("user_id", itoa(int(msg.UserID)))
	for _, u := range msg.UserIDs {
		args.Add("user_ids", itoa(int(u)))
	}
	args.Set("peer_id", itoa(int(msg.PeerID)))
	for _, p := range msg.PeerIDs {
		args.Add("peer_ids", itoa(int(p)))
	}
	args.Set("domain", msg.Domain)
	args.Set("chat_id", itoa(int(msg.ChatID)))

	args.Set("message", msg.Message)
	args.Set("lat", itoa(msg.Lat))
	args.Set("long", itoa(msg.Long))
	args.Set("attachment", msg.Attachment...)
	args.Set("sticker_id", itoa(msg.StickerID))
	args.Set("keyboard", msg.Keyboard.String())

	args.Set("reply_to", itoa(msg.ReplyTo))
	for _, f := range msg.ForwardMessages {
		args.Add("peer_ids", itoa(f))
	}

	args.Set("dont_parse_links", msg.DontParseLinks.String())
	args.Set("disable_mentions", msg.DisableMentions.String())

	var id int
	err := m.client.Method("messages.send", args, &id)
	return id, err
}

// Members struct.
type Members struct {
	Items    []Member `json:"items"`
	Profiles []User   `json:"profiles"`
	Groups   []Group  `json:"groups"`
	Count    int      `json:"count"`
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
func (m *MessagesProvider) GetConversationMembers(peerID ID) (*Members, error) {
	args := NewArgs().Set("peer_id", itoa(int(peerID)))
	var mem Members
	err := m.client.Method("messages.getConversationMembers", args, &m)
	return &mem, err
}

// RemoveChatUser allows the current user to leave a chat or, if the current user started the chat,
// allows the user to remove another user from the chat.
func (m *MessagesProvider) RemoveChatUser(chatID ChatID, memberID ID) error {
	args := NewArgs()
	args.Set("chat_id", itoa(int(chatID)))
	args.Set("member_id", itoa(int(memberID)))
	err := m.client.Method("messages.removeChatUser", args, nil)
	return err
}
