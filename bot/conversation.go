package bot

import (
	"bytes"
	"io"
	"io/fs"

	"github.com/Toffee-iZt/HwBot/vkapi"
)

// Message struct.
type Message struct {
	Text        string
	Lat, Long   float64
	Attachments []*Attachment
	Keyboard    *Keyboard
	ReplyTo     int
	Forward     []int
	Mentions    bool
}

// NewAttachment makes new attachment.
func NewAttachment(name string, reader io.Reader) *Attachment {
	if reader == nil {
		return nil
	}
	return &Attachment{
		name:   name,
		reader: reader,
	}
}

// NewAttachmentFile makes new attachment from file.
func NewAttachmentFile(file fs.File) (*Attachment, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return &Attachment{
		name:   stat.Name(),
		reader: bytes.NewReader(data),
	}, nil
}

// Attachment struct.
type Attachment struct {
	name   string
	reader io.Reader
}

// Conversation struct.
type Conversation struct {
	peer vkapi.ID
	api  *vkapi.Client
}

// GetMembers returns a list of conversation users.
func (c *Conversation) GetMembers() (*vkapi.Members, *vkapi.Error) {
	return c.api.GetConversationMembers(c.peer)
}

func (c *Conversation) prepareMessage(msg Message) (vkapi.OutMessageContent, *vkapi.Error) {
	var atts []string
	if len(msg.Attachments) > 0 {
		atts = make([]string, len(msg.Attachments))
		for i, a := range msg.Attachments {
			attstr, vkerr := c.api.UploadMessagesPhoto(c.peer, a.name, a.reader)
			if vkerr != nil {
				return vkapi.OutMessageContent{}, vkerr
			}
			atts[i] = attstr
		}
	}

	var kb vkapi.JSONData
	if msg.Keyboard != nil {
		kb = msg.Keyboard.Build()
	}

	return vkapi.OutMessageContent{
		Message:         msg.Text,
		Lat:             msg.Lat,
		Long:            msg.Long,
		Attachment:      atts,
		Keyboard:        kb,
		ReplyTo:         msg.ReplyTo,
		ForwardMessages: msg.Forward,
		DontParseLinks:  true,
		DisableMentions: !msg.Mentions,
	}, nil
}

// SendMessage sends a message.
func (c *Conversation) SendMessage(msg Message) *vkapi.Error {
	m, vkerr := c.prepareMessage(msg)
	if vkerr != nil {
		return vkerr
	}

	return c.api.SendMessage(vkapi.OutMessage{
		PeerID:            c.peer,
		OutMessageContent: m,
	})
}

// EditMessage edits a message.
func (c *Conversation) EditMessage(convMessageID int, msg Message) (bool, *vkapi.Error) {
	m, vkerr := c.prepareMessage(msg)
	if vkerr != nil {
		return false, vkerr
	}
	return c.api.EditMessage(c.peer, convMessageID, m)
}
