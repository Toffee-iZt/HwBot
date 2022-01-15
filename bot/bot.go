package bot

import (
	"context"
	"errors"
	"strings"

	"github.com/Toffee-iZt/HwBot/common"
	"github.com/Toffee-iZt/HwBot/logger"
	"github.com/Toffee-iZt/HwBot/vkapi"
	"github.com/Toffee-iZt/HwBot/vkapi/longpoll"
)

// New creates new bot.
func New(vk *vkapi.Client, l *logger.Logger) *Bot {
	return &Bot{
		log:      l,
		vk:       vk,
		lp:       longpoll.New(vk),
		commands: make(map[string]*command),
	}
}

// Bot is vk bot.
type Bot struct {
	sync common.Sync

	log    *logger.Logger
	modlog *logger.Logger
	vk     *vkapi.Client
	lp     *longpoll.LongPoll

	modules []*Module

	commands map[string]*command
}

type command struct {
	*Command
	*Module
}

// Module struct.
type Module struct {
	Name      string
	Init      func(*logger.Logger) bool
	Terminate func()
	Commands  []*Command
	log       *logger.Logger
}

// Command respresents conversation command.
type Command struct {
	Run         func(*Context, *NewMessage, []string)
	Cmd         []string
	Description string
	Help        string
	InPrivate   bool
	InChat      bool
	NoPrefix    bool
}

func (b *Bot) close(err error) {
	b.sync.LockClose(func() error {
		b.commands = nil
		for i := range b.modules {
			t := b.modules[i].Terminate
			if t != nil {
				t()
			}
		}
		b.modules = nil
		return err
	})
}

// Done returns channel that is closed whed bot stopped.
func (b *Bot) Done() <-chan struct{} {
	return b.sync.Done()
}

// Err returns bot error.
func (b *Bot) Err() error {
	return b.sync.Err()
}

// Run starts bot.
func (b *Bot) Run(ctx context.Context, sync bool, mods ...*Module) error {
	if !b.sync.Init() {
		return errors.New("bot is already running")
	}

	b.modules = make([]*Module, 0, len(mods))

	b.log.Info("loading modules")
	b.modlog = b.log.Child("MODULES")
	for _, m := range mods {
		m.log = b.modlog.Child(m.Name)
		if m.Init != nil && !m.Init(m.log) {
			b.modlog.Warn(`module "%s" init failed`, m.Name)
			continue
		}

		for _, cmd := range m.Commands {
			c := &command{
				Command: cmd,
				Module:  m,
			}
			for _, alias := range cmd.Cmd {
				b.commands[alias] = c
			}
		}
		b.modules = append(b.modules, m)
		b.modlog.Info(`module "%s" inited`, m.Name)
	}

	if sync {
		b.run(ctx)
		return b.Err()
	}
	go b.run(ctx)
	return nil
}

func (b *Bot) run(ctx context.Context) {
	ch := b.lp.Run(ctx, 25)
	for {
		select {
		case e := <-ch:
			go b.handle(ctx, e)
		case <-b.lp.Done():
			b.close(b.lp.Err())
			return
		}
	}
}

func (b *Bot) handle(ctx context.Context, e longpoll.Event) {
	switch e.Type {
	case longpoll.TypeMessageEvent:
		// TODO
		//cbNew := e.Object.(*longpoll.MessageEvent)
		//p, err := events.ParsePayload(cbNew.Payload)
		//if err != nil || p == nil || p.Module == "" {
		//	b.log.Warn("parsing message payload: empty or invalid format\n%d in %d\n%s", cbNew.UserID, cbNew.PeerID, string(cbNew.Payload))
		//}
		//mod, ok := b.modules[p.Module]
		//if !ok {
		//	b.log.Warn("parsing message payload: invalid module name\n%d in %d\n%s", cbNew.UserID, cbNew.PeerID, string(cbNew.Payload))
		//}
		//mod.Process()
	case longpoll.TypeMessageNew:
		msgNew := e.Object.(*longpoll.MessageNew)

		cmd, args := b.getCommand(&msgNew.Message)
		if cmd == nil {
			return
		}

		bctx := makeContext(b, msgNew.Message.PeerID, cmd, msgNew)

		id := int(msgNew.Message.PeerID)
		if id > 2e9 && !cmd.InChat || id < 2e9 && !cmd.InPrivate {
			bctx.ReplyText("Команда недоступна в данном типе чата")
		}

		cmd.Run(bctx, msgNew, args)
	}
}

// CommandPrefixes are the characters with which commands must begin.
const CommandPrefixes = "/!"

func (b *Bot) getCommand(msg *longpoll.Message) (*command, []string) {
	if n := len(msg.Text); n > 1 && strings.IndexByte(CommandPrefixes, msg.Text[0]) != -1 {
		if n > 300 {
			msg.Text = msg.Text[:300]
		}
		s := strings.Split(msg.Text, " ")
		return b.commands[s[0][1:]], s[1:]
	}
	return nil, nil
}

// ModList returns loaded modules list.
func (b *Bot) ModList() []*Module {
	return b.modules
}

// API returns vk api instance.
func (b *Bot) API() *vkapi.Client {
	return b.vk
}
