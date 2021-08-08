package bot

import (
	"HwBot/bot/api"
	"HwBot/vkapi"
	"HwBot/vkapi/longpoll"
	"context"
	"fmt"
	"strings"
)

func (b *Bot) handle(ctx context.Context, e longpoll.Event) {
	switch e.Type {
	case longpoll.TypeMessageNew:
		b.onMessage(ctx, &api.IncomingMessage{e.Object.(*longpoll.MessageNew)})
	case longpoll.TypeMessageEvent:
		b.onCallback(ctx, &api.CallbackMessage{e.Object.(*longpoll.MessageEvent)})
	}
}

func (b *Bot) execCommand(msg *api.IncomingMessage) {
	s := strings.Split(msg.Message.Text, " ")
	cmd := s[0][1:]
	if cmd == "help" {
		b.sendHelp(msg, s[1:])
		return
	}
	c := b.commands[cmd]
	if c == nil {
		return
	}

	c.Run(b, msg, s[1:])
}

func (b *Bot) onMessage(ctx context.Context, msg *api.IncomingMessage) {
	m := &msg.Message
	if n := len(m.Text); n > 1 && strings.IndexByte(b.cfg.Prefixes, m.Text[0]) != -1 {
		if n > 300 {
			m.Text = m.Text[:300]
		}
		go b.execCommand(msg)
		return
	}

	for i := range b.modules {
		m := b.modules[i]
		if m.OnMessage != nil {
			m.OnMessage(b, msg)
		}
	}
}

func (b *Bot) onCallback(ctx context.Context, cb *api.CallbackMessage) {}

func (b *Bot) sendHelp(msg *api.IncomingMessage, a []string) {
	if len(a) > 0 {
		c := b.commands[a[0]]
		if c != nil {
			b.API().Messages.Send(vkapi.MessagePeer{PeerID: msg.Message.PeerID}, vkapi.MessageContent{Message: fmt.Sprintf("%s - %s\n%s", c.Cmd, c.Desc, c.Help)})
			return
		}
	}
	var h string

	for _, c := range b.commands {
		h += fmt.Sprintf("%s - %s\n", c.Cmd, c.Desc)
	}
	b.API().Messages.Send(vkapi.MessagePeer{PeerID: msg.Message.PeerID}, vkapi.MessageContent{Message: h})
}
