package bot

import (
	"context"
	"strings"

	"github.com/Toffee-iZt/HwBot/vkapi/longpoll"
)

// CommandPrefixes are the characters with which commands must begin.
const CommandPrefixes = "/!"

// IncomingMessage event.
type IncomingMessage struct {
	*longpoll.MessageNew
}

// CallbackMessage event.
type CallbackMessage struct {
	*longpoll.MessageEvent
}

func (b *Bot) handle(ctx context.Context, e longpoll.Event) {
	switch e.Type {
	case longpoll.TypeMessageNew:
		b.onMessage(ctx, &IncomingMessage{e.Object.(*longpoll.MessageNew)})
	case longpoll.TypeMessageEvent:
		b.onCallback(ctx, &CallbackMessage{e.Object.(*longpoll.MessageEvent)})
	}
}

func (b *Bot) execCommand(msg *IncomingMessage) {
	s := strings.Split(msg.Message.Text, " ")
	c := b.commands[s[0][1:]]
	if isuser := msg.Message.PeerID < 2e9; c == nil || !c.Priv && isuser || !c.Chat && !isuser {
		return
	}

	c.Run(b, msg, s[1:])
}

func (b *Bot) onMessage(ctx context.Context, msg *IncomingMessage) {
	m := &msg.Message
	if n := len(m.Text); n > 1 && strings.IndexByte(CommandPrefixes, m.Text[0]) != -1 {
		if n > 300 {
			m.Text = m.Text[:300]
		}
		go b.execCommand(msg)
		return
	}

	// debug
	u, err := b.API().Users.Get([]int{m.FromID})
	if err != nil {
		b.log.Error(err.Error())
		return
	}

	b.log.Info(":::DEBUG:::New message: from: %s %s\n%s", u[0].FirstName, u[0].LastName, m.Text)
}

func (b *Bot) onCallback(ctx context.Context, cb *CallbackMessage) {}
