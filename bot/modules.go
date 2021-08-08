package bot

import (
	"HwBot/bot/api"
	"HwBot/common/rt"
	"HwBot/logger"
	"HwBot/vkapi"
)

// API returns vk client instance.
func (b *Bot) API() *vkapi.Client {
	return b.vk
}

// SimpleReply ...
func (b *Bot) SimpleReply(m *api.IncomingMessage, text string) {
	_, vkerr := b.API().Messages.Send(
		vkapi.MessagePeer{
			PeerID: m.Message.PeerID,
		},
		vkapi.MessageContent{
			Message: text,
		},
	)
	if vkerr != nil {
		c := rt.Caller()
		b.log.Error("simple reply error: %s %s", vkerr.Error(), c.Function)
	}
}

// NewCallback ...
func (b *Bot) NewCallback() {}

func (b *Bot) initModule(m *api.Module, l *logger.Logger) {
	if m.Init != nil {
		ok := m.Init(b, l.Child(m.Name))
		if !ok {
			b.log.Warn(`module "%s" init failed`, m.Name)
			return
		}
	}

	for i := range m.Commands {
		b.regCommand(m.Commands[i])
	}
	b.modules = append(b.modules, m)
}

func (b *Bot) regCommand(cmd *api.Command) {
	b.commands[cmd.Cmd] = cmd
}
