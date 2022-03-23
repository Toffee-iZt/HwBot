package bot

import (
	"context"
	"errors"

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
		commands: make(map[string]*Command),
	}
}

// Bot is vk bot.
type Bot struct {
	sync common.Sync

	log *logger.Logger
	vk  *vkapi.Client
	lp  *longpoll.LongPoll

	modules []*Module

	commands map[string]*Command
}

// Module struct.
type Module struct {
	Name      string
	Init      func(*logger.Logger) bool
	Terminate func()
	Callback  func(ctx *CallbackContext, user vkapi.UserID, msg int, payload vkapi.JSONData)
	Commands  []*Command
	log       *logger.Logger
}

// Command respresents conversation command.
type Command struct {
	Run         func(*MessageContext, *NewMessage, []string)
	Cmd         []string
	Description string
	Help        string
	InPrivate   bool
	InChat      bool
	NoPrefix    bool
	module      *Module
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
	modlog := b.log.Child("MODULES")
	for _, m := range mods {
		m.log = modlog.Child(m.Name)
		if m.Init != nil && !m.Init(m.log) {
			m.log.Warn("init failed")
			continue
		}

		for _, cmd := range m.Commands {
			cmd.module = m
			for _, alias := range cmd.Cmd {
				b.commands[alias] = cmd
			}
		}
		b.modules = append(b.modules, m)
		m.log.Info("init successful")
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
			go func() {
				switch e.Type {
				case longpoll.EventTypeMessageEvent:
					b.onCallback(e.Object.(*longpoll.MessageEvent))
				case longpoll.EventTypeMessageNew:
					b.onMessage(e.Object.(*longpoll.MessageNew))
				}
			}()
		case <-b.lp.Done():
			b.close(b.lp.Err())
			return
		}
	}
}

// ModList returns loaded modules list.
func (b *Bot) ModList() []*Module {
	return b.modules
}
