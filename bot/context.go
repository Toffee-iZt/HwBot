package bot

import (
	"fmt"
	"runtime"

	"github.com/Toffee-iZt/HwBot/common/rt"
	"github.com/Toffee-iZt/HwBot/logger"
	"github.com/Toffee-iZt/HwBot/vkapi"
)

type eventctx struct {
	Conv *Conversation
	mod  *Module
}

func (c *eventctx) close() {
	c.Conv = nil
	c.mod = nil
	runtime.Goexit()
}

func (c *eventctx) errlog(f string, vkerr *vkapi.Error, a ...interface{}) {
	if vkerr != nil {
		c.mod.log.Error(f, vkerr.Error(), fmt.Sprint(a...))
	}
}

// Log returns module logger.
func (c *eventctx) Log() *logger.Logger {
	return c.mod.log
}

// API returns vk client.
func (c *eventctx) API() *vkapi.Client {
	return c.Conv.api
}

// KeyboardGenerator returns empty keyboard generator.
func (c *eventctx) KeyboardGenerator(oneTime bool, inline bool) *Keyboard {
	kb := vkapi.NewKeyboard(oneTime, inline)
	return &Keyboard{
		mod: c.mod,
		kb:  kb,
	}
}

// Reply replies with a message and closes context.
func (c *eventctx) Reply(msg Message) {
	vkerr := c.Conv.SendMessage(msg)
	c.errlog("context reply error: %s %s", vkerr, rt.Caller().Function)
	c.close()
}

// ReplyMessage replies with a text and closes context.
func (c *eventctx) ReplyMessage(text string, attachments ...*Attachment) {
	vkerr := c.Conv.SendMessage(Message{
		Text:        text,
		Attachments: attachments,
	})
	c.errlog("context reply error: %s %s", vkerr, rt.Caller().Function)
	c.close()
}

// Context struct.
type Context struct {
	eventctx
}

// ReplyText replies to a message with a text and closes context.
func (c *Context) ReplyText(text string) {
	c.ReplyMessage(text)
}

// ReplyError send error as vk message and closes context.
func (c *Context) ReplyError(f string, a ...interface{}) {
	c.ReplyText(fmt.Sprintf(f, a...))
}

// CallbackContext struct.
type CallbackContext struct {
	eventctx

	eventID string
	userID  vkapi.UserID
	msgID   int
}

// ReplyCallback replies to a callback event and closes context.
func (c *CallbackContext) ReplyCallback(data *vkapi.EventData) {
	vkerr := c.Conv.api.SendMessageEventAnswer(c.eventID, c.userID, c.Conv.peer, data)
	c.errlog("callback context reply error: %s %s", vkerr, rt.Caller().Function)
	c.close()
}

// MessageID returns message id.
func (c *CallbackContext) MessageID() int {
	return c.msgID
}

// UserID returns user id.
func (c *CallbackContext) UserID(data *vkapi.EventData) vkapi.UserID {
	return c.userID
}
