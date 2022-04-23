package conversation

import (
	"bytes"
	"io"
	"io/fs"

	"github.com/Toffee-iZt/HwBot/vkapi"
	"github.com/Toffee-iZt/HwBot/vkapi/keyboard"
	"github.com/Toffee-iZt/HwBot/vkapi/upload"
)

// Message struct.
type Message struct {
	Text     string
	Photos   []*Photo
	Keyboard *keyboard.Keyboard
	ReplyTo  int
	Forward  []int
	Mentions bool
}

// Photo type
type Photo struct {
	fname string
	data  io.Reader
}

// NewPhoto creates photo object.
func NewPhoto(fname string, data io.Reader) *Photo {
	if data == nil {
		return nil
	}
	return &Photo{
		fname: fname,
		data:  data,
	}
}

// NewPhotoFile creates photo object from file.
func NewPhotoFile(file fs.File) (*Photo, error) {
	if file == nil {
		return nil, nil
	}
	fi, err := file.Stat()
	if err != nil || fi.Size() == 0 {
		return nil, err
	}
	buf := bytes.NewBuffer(make([]byte, 0, fi.Size()))
	_, err = buf.ReadFrom(file)
	if err != nil {
		return nil, err
	}
	return &Photo{
		fname: fi.Name(),
		data:  buf,
	}, nil
}

// New creates new chat object.
func New(peer vkapi.ID, api *vkapi.Client) *Conversation {
	return &Conversation{
		peer: peer,
		api:  api,
	}
}

// Conversation struct.
type Conversation struct {
	peer vkapi.ID
	api  *vkapi.Client
}

// GetMembers returns a list of chat users.
func (c *Conversation) GetMembers() (*vkapi.Members, *vkapi.Error) {
	return c.api.GetConversationMembers(c.peer)
}

func (c *Conversation) prepareMessage(msg Message) (vkapi.OutMessageContent, *vkapi.Error) {
	atts := make([]string, len(msg.Photos))
	for i, a := range msg.Photos {
		vkstr, vkerr := upload.MessagesPhoto(c.api, c.peer, a.fname, a.data)
		if vkerr != nil {
			return vkapi.OutMessageContent{}, vkerr
		}
		atts[i] = vkstr
	}

	var kb vkapi.JSONData
	if msg.Keyboard != nil {
		kb = msg.Keyboard.Data()
	}

	return vkapi.OutMessageContent{
		Message:         msg.Text,
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
