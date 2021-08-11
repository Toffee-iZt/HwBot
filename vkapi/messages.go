package vkapi

import (
	"strconv"
	"sync/atomic"

	"github.com/Toffee-iZt/HwBot/vkapi/vktypes/keyboard"
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

// MessagePeer struct.
type MessagePeer struct {
	UserID  int
	UserIDs []int
	PeerID  int
	PeerIDs []int
	Domain  string
	ChatID  int
}

// MessageContent struct.
type MessageContent struct {
	Message    string
	Lat, Long  int
	Attachment []string
	StickerID  int
	Keyboard   *keyboard.Keyboard
	//Template map[string]interface{}
	//Payload       map[string]interface{} // ???
	//ContentSource map[string]interface{} // ???

	Forward struct {
		ReplyTo         int
		ForwardMessages []int
		//Forward common.JSONData
	}

	Meta struct {
		DontParseLinks  baseBoolInt
		DisableMentions baseBoolInt

		//GroupID         int

		//Intent      string // ???
		//SubscribeID int   // ???
	}
}

// Send sends a message.
func (m *MessagesProvider) Send(peer MessagePeer, msg MessageContent) (int, error) {
	if msg.Message == "" && msg.Attachment == nil {
		return 0, nil
	}

	if len(msg.Attachment) > 10 {
		msg.Attachment = msg.Attachment[:10]
	}

	args := NewArgs()
	argsSetAny(args, "random_id", int(atomic.AddInt32(&m.randomID, 1)))

	argsSetAny(args, "user_id", peer.UserID)
	argsSetAny(args, "peer_id", peer.PeerID)
	argsSetAny(args, "peer_ids", peer.PeerIDs)
	argsSetAny(args, "domain", peer.Domain)
	argsSetAny(args, "chat_id", peer.ChatID)
	argsSetAny(args, "user_ids", peer.UserIDs)

	argsSetAny(args, "message", msg.Message)
	argsSetAny(args, "lat", msg.Lat)
	argsSetAny(args, "long", msg.Long)
	argsSetAny(args, "attachment", msg.Attachment)
	argsSetAny(args, "sticker_id", msg.StickerID)
	argsSetAny(args, "keyboard", msg.Keyboard)

	argsSetAny(args, "reply_to", msg.Forward.ReplyTo)
	argsSetAny(args, "forward_messages", msg.Forward.ForwardMessages)

	argsSetAny(args, "dont_parse_links", msg.Meta.DontParseLinks)
	argsSetAny(args, "disable_mentions", msg.Meta.DisableMentions)

	var id int
	err := m.client.Method("messages.send", args, &id)
	return id, err
}

// Members ...
type Members struct {
	Items    []Member `json:"items"`
	Profiles []User   `json:"profiles"`
	Groups   []Group  `json:"groups"`
	Count    int      `json:"count"`
}

// Member ...
type Member struct {
	MemberID  int   `json:"member_id"`
	InvitedBy int   `json:"invited_by"`
	JoinDate  int64 `json:"join_date"`
	IsAdmin   bool  `json:"is_admin"`
	IsOwner   bool  `json:"is_owner"`
	CanKick   bool  `json:"can_kick"`
}

// GetChatMembers ...
func (m *MessagesProvider) GetChatMembers(chatID int) (*Members, error) {
	args := NewArgs().Set("peer_id", strconv.Itoa(chatID+2e9))
	var mem Members
	err := m.client.Method("messages.getConversationMembers", args, &m)
	return &mem, err
}

// Kick ...
func (m *MessagesProvider) Kick(chatID int, memberID int) error {
	args := NewArgs()
	args.Set("chat_id", strconv.Itoa(chatID))
	args.Set("member_id", strconv.Itoa(memberID))
	err := m.client.Method("messages.removeChatUser", args, nil)
	return err
}
