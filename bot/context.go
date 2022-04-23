package bot

import (
	"runtime"

	"github.com/Toffee-iZt/HwBot/bot/conversation"
	"github.com/Toffee-iZt/HwBot/common/rt"
	"github.com/Toffee-iZt/HwBot/logger"
	"github.com/Toffee-iZt/HwBot/vkapi"
	"github.com/Toffee-iZt/HwBot/vkapi/longpoll"
)

type eventctx struct {
	con *conversation.Conversation
	mod *Module
	api *vkapi.Client
}

func (c *eventctx) errlog(vkerr *vkapi.Error) {
	if vkerr != nil {
		c.mod.log.Error("context error: %d %s\n%s %s\n%s",
			vkerr.Code,
			vkerr.Message,
			vkerr.Method,
			vkerr.Args,
			rt.Caller(1).Function,
		)
	}
}

// Close closes current context.
func (c *eventctx) Close() {
	runtime.Goexit()
}

// Log returns logger for module.
func (c *eventctx) Log() *logger.Logger {
	return c.mod.log
}

// API returns vk client.
func (c *eventctx) API() *vkapi.Client {
	return c.api
}

// Conversation returns chat object for current chat.
func (c *eventctx) Conversation() *conversation.Conversation {
	return c.con
}

// Payload wraps payload.
func (c *eventctx) Payload(i interface{}) vkapi.JSONData {
	return wrap(c.mod, i)
}

func newContext(api *vkapi.Client, mod *Module, peer vkapi.ID) *Context {
	return &Context{
		eventctx: eventctx{
			mod: mod,
			api: api,
			con: conversation.New(peer, api),
		},
	}
}

// Context type.
type Context struct {
	eventctx
}

// Reply replies with a message and closes context.
func (c *Context) Reply(msg conversation.Message) {
	vkerr := c.con.SendMessage(msg)
	c.errlog(vkerr)
	c.Close()
}

// ReplyText replies to a message with a text and closes context.
func (c *Context) ReplyText(text string) {
	c.Reply(conversation.Message{Text: text})
}

func newCallback(api *vkapi.Client, mod *Module, event *longpoll.MessageEvent) *Callback {
	return &Callback{
		eventctx: eventctx{
			mod: mod,
			api: api,
			con: conversation.New(event.PeerID, api),
		},
		event: event,
	}
}

// Callback type.
type Callback struct {
	eventctx
	event *longpoll.MessageEvent
}

// Close replies without action to a callback event and closes context.
func (c *Callback) Close() {
	c.reply(nil)
}

// Edit edits a message with keyboard and closes context.
func (c *Callback) Edit(msg conversation.Message) {
	if c.event.ConversationMessageID == 0 {
		c.mod.log.Warn("callback context edit message: keyboard is not inlined")
		c.Close()
	}
	_, vkerr := c.con.EditMessage(c.event.ConversationMessageID, msg)
	c.errlog(vkerr)
	c.Close()
}

// ShowSnackbar replies to a callback event and closes context.
func (c *Callback) ShowSnackbar(text string) {
	c.reply(&vkapi.EventData{
		Type: vkapi.EventDataTypeShowSnackbar,
		Text: text,
	})
}

// OpenLink replies to a callback event and closes context.
func (c *Callback) OpenLink(link string) {
	c.reply(&vkapi.EventData{
		Type: vkapi.EventDataTypeOpenLink,
		Link: link,
	})
}

func (c *Callback) reply(data *vkapi.EventData) {
	vkerr := c.api.SendMessageEventAnswer(c.event.EventID, c.event.UserID, c.event.PeerID, data)
	c.errlog(vkerr)
	c.Close()
}

// MessageID returns message id.
func (c *Callback) MessageID() int {
	return c.event.ConversationMessageID
}

// UserID returns user id.
func (c *Callback) UserID() vkapi.UserID {
	return c.event.UserID
}
