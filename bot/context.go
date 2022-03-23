package bot

import (
	"fmt"
	"runtime"

	"github.com/Toffee-iZt/HwBot/common/rt"
	"github.com/Toffee-iZt/HwBot/vkapi"
)

type eventctx struct {
	Conv *Conversation
	bot  *Bot
	mod  *Module
}

func (c *eventctx) close() {
	c.Conv = nil
	c.bot = nil
	c.mod = nil
	runtime.Goexit()
}

func (c *eventctx) errlog(f string, vkerr *vkapi.Error, a ...interface{}) {
	if vkerr != nil {
		c.mod.log.Error(f, vkerr.Error(), fmt.Sprint(a...))
	}
}

// Bot returns bot instance.
func (c *eventctx) Bot() *Bot {
	return c.bot
}

// API returns vk client.
func (c *eventctx) API() *vkapi.Client {
	return c.bot.vk
}

// KeyboardGenerator returns empty keyboard generator.
func (c *eventctx) KeyboardGenerator(oneTime bool, inline bool) *Keyboard {
	kb := vkapi.NewKeyboard(oneTime, inline)
	kb.NextRow()
	return &Keyboard{
		mod: c.mod,
		kb:  kb,
	}
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

// MessageContext struct.
type MessageContext struct {
	eventctx
}

// ReplyText replies to a message with a text and closes context.
func (c *MessageContext) ReplyText(text string) {
	c.ReplyMessage(text)
}

// ReplyError send error as vk message and closes context.
func (c *MessageContext) ReplyError(f string, a ...interface{}) {
	c.ReplyText(fmt.Sprintf(f, a...))
}

// CallbackContext struct.
type CallbackContext struct {
	eventctx

	eventID string
	userID  vkapi.UserID
}

// ReplyCallback replies to a callback event and closes context.
func (c *CallbackContext) ReplyCallback(data *vkapi.EventData) {
	vkerr := c.Conv.api.SendMessageEventAnswer(c.eventID, c.userID, c.Conv.peer, data)
	c.errlog("callback context reply error: %s %s", vkerr, rt.Caller().Function)
	c.close()
}
