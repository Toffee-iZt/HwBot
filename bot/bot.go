package bot

import (
	"context"
	"strings"

	"github.com/Toffee-iZt/HwBot/common/strbytes"
	"github.com/Toffee-iZt/HwBot/logger"
	"github.com/Toffee-iZt/HwBot/vkapi"
	"github.com/Toffee-iZt/HwBot/vkapi/longpoll"
)

// New creates new bot.
func New(vk *vkapi.Client, l *logger.Logger) *Bot {
	return &Bot{
		log: l,
		vk:  vk,
	}
}

// Bot is vk bot.
type Bot struct {
	log     *logger.Logger
	vk      *vkapi.Client
	modules []*Module
}

func (b *Bot) close() {
	for _, m := range b.modules {
		if m.Terminate != nil {
			m.Terminate()
		}
	}
	b.modules = nil
}

// Run starts bot.
func (b *Bot) Run(ctx context.Context, mods ...*Module) (err error) {
	b.modules = make([]*Module, 0, len(mods))
	b.log.Info("loading modules")
	for _, m := range mods {
		if m.Init != nil {
			if err := m.Init(); err != nil {
				b.log.Error("module '%s' init failed: %s", m.Name, err.Error())
				return nil
			}
		}
		m.log = b.log.Child(m.Name)
		b.modules = append(b.modules, m)
		b.log.Info("module '%s' init successful", m.Name)
	}
	err = b.run(ctx)
	return
}

func (b *Bot) run(ctx context.Context) error {
	ch := longpoll.Run(ctx, b.vk, 25)
	for {
		ev, ok := <-ch
		if !ok {
			// context is done
			b.close()
			return context.Canceled
		}
		go func() {
			switch ev.Type {
			case longpoll.EventTypeMessageEvent:
				b.onCallback(ev.Object.(*longpoll.MessageEvent))
			case longpoll.EventTypeMessageNew:
				b.onMessage(ev.Object.(*longpoll.MessageNew))
			}
		}()
	}
}

func (b *Bot) find(cmd string) (m *Module, c *Command) {
	for _, m = range b.modules {
		for _, c = range m.Commands {
			if strbytes.Has(c.Cmd, cmd) {
				return
			}
		}
	}
	return
}

func (b *Bot) onMessage(msg *longpoll.MessageNew) {
	// MessageLimit is a limit to the number of characters of the input message.
	const MessageLimit = 300

	n := len(msg.Message.Text)
	if n == 0 || strings.IndexByte(Prefixes, msg.Message.Text[0]) == -1 {
		return
	}
	if n > MessageLimit {
		msg.Message.Text = msg.Message.Text[:MessageLimit]
	}

	split := strings.Split(msg.Message.Text[1:], " ")

	mod, cmd := b.find(split[0])
	if cmd == nil {
		return
	}
	message := Message{
		Message: msg.Message,
		Args:    split[1:],
	}

	c := newContext(b.vk, mod, msg.Message.PeerID)
	id := message.PeerID
	if id.ToChat() != 0 && !cmd.Options.Has(OptionInChat) ||
		id.ToUser() != 0 && !cmd.Options.Has(OptionInDialog) {
		c.ReplyText("Команда недоступна в данном типе чата")
	}

	cmd.Run(c, &message)
}

func (b *Bot) onCallback(event *longpoll.MessageEvent) {
	modname, data := unwrap(event.Payload)
	if modname == "" {
		return
	}
	for _, mod := range b.modules {
		if mod.Name == modname {
			if mod.Callback == nil {
				return
			}
			mod.Callback(newCallback(b.vk, mod, event), data)
			break
		}
	}
}

// ModList returns loaded modules list.
func (b *Bot) ModList() []*Module {
	return b.modules
}
