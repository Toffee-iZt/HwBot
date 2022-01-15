package bot

import (
	"fmt"
	"io"
	"runtime"

	"github.com/Toffee-iZt/HwBot/common/rt"
	"github.com/Toffee-iZt/HwBot/logger"
	"github.com/Toffee-iZt/HwBot/vkapi"
)

func makeContext(b *Bot, peerID vkapi.ID, cmd *command) *Context {
	return &Context{
		peer: peerID,
		bot:  b,
		log:  cmd.log,
		cmd:  cmd.Command,
	}
}

// Context struct.
type Context struct {
	peer vkapi.ID
	bot  *Bot
	log  *logger.Logger
	cmd  *Command
}

func (c *Context) close() {
	c.peer = 0
	c.bot = nil
	c.cmd = nil
	runtime.Goexit()
}

func (c *Context) errlog(fmt string, err error, a ...interface{}) {
	if err != nil {
		c.log.Error(fmt, a...)
	}
}

// Reply replies to a message with a text and closes bridge.
func (c *Context) Reply(text string, attachments ...string) {
	_, vkerr := c.SendMessage(vkapi.OutMessageContent{
		Message:    text,
		Attachment: attachments,
	})
	c.errlog("context reply error: %s %s", vkerr, rt.Caller().Function)
	c.close()
}

// ReplyText replies to a message with a text and closes bridge.
func (c *Context) ReplyText(text string) {
	c.Reply(text)
}

// ReplyError send error as vk message and closes context.
func (c *Context) ReplyError(f string, a ...interface{}) {
	c.ReplyText(fmt.Sprintf(f, a...))
}

// BotInstance returns bot instance.
func (c *Context) BotInstance() *Bot {
	return c.bot
}

// UploadAttachment uploads attachment and returns vk attachment string.
func (c *Context) UploadAttachment(name string, data io.Reader) (string, error) {
	return c.bot.vk.UploadMessagesPhoto(c.peer, name, data)
}

// GetMembers returns a list of conversation users.
func (c *Context) GetMembers() (*vkapi.Members, error) {
	return c.bot.vk.GetConversationMembers(c.peer)
}

// SendMessage sends a message.
func (c *Context) SendMessage(msg vkapi.OutMessageContent) (int, error) {
	return c.bot.vk.Send(vkapi.OutMessage{
		PeerID:            c.peer,
		OutMessageContent: msg,
	})
}
