package bot

import (
	"github.com/Toffee-iZt/HwBot/common/rt"
	"github.com/Toffee-iZt/HwBot/logger"
	"github.com/Toffee-iZt/HwBot/vkapi"
)

// Module struct.
type Module struct {
	Name      string
	Init      func(*Bot, *logger.Logger) bool
	Terminate func()
	Commands  []*Command
}

// Command respresents conversation command.
type Command struct {
	Run  func(*Bot, *IncomingMessage, []string)
	Cmd  string
	Desc string
	Help string
	Priv bool
	Chat bool
}

// API returns vk client instance.
func (b *Bot) API() *vkapi.Client {
	return b.vk
}

// Modules returns modules list.
func (b *Bot) Modules() []*Module {
	return b.modules
}

// SimpleReply ...
func (b *Bot) SimpleReply(m *IncomingMessage, text string) {
	_, vkerr := b.API().Messages.Send(vkapi.OutMessage{
		PeerID:  m.Message.PeerID,
		Message: text,
	})
	if vkerr != nil {
		c := rt.Caller()
		b.log.Error("simple reply error: %s %s", vkerr.Error(), c.Function)
	}
}

// NewCallback ...
func (b *Bot) NewCallback() {}

func (b *Bot) initModule(m *Module, l *logger.Logger) {
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
	b.log.Info(`module "%s" inited`, m.Name)
}

func (b *Bot) regCommand(cmd *Command) {
	b.commands[cmd.Cmd] = cmd
}
