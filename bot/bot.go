package bot

import (
	"context"

	"github.com/Toffee-iZt/HwBot/common"
	"github.com/Toffee-iZt/HwBot/logger"
	"github.com/Toffee-iZt/HwBot/vkapi"
	"github.com/Toffee-iZt/HwBot/vkapi/longpoll"
)

// New creates new bot.
func New(vk *vkapi.Client, pref []byte, w *logger.Writer) *Bot {
	return &Bot{
		log:      logger.New(w, "BOT"),
		pref:     pref,
		vk:       vk,
		lp:       longpoll.New(vk, 25),
		commands: make(map[string]*Command),
	}
}

// Bot is vk bot.
type Bot struct {
	sync common.Sync

	pref []byte

	log *logger.Logger
	vk  *vkapi.Client
	lp  *longpoll.LongPoll

	commands map[string]*Command
	modules  []*Module
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
func (b *Bot) Run(ctx context.Context, mods ...*Module) bool {
	if !b.sync.Init() {
		return false
	}

	b.log.Info("loading modules")
	modlog := b.log.Child("MODULES")
	for _, m := range mods {
		b.initModule(m, modlog)
	}

	go b.run(ctx)
	return true
}

func (b *Bot) run(ctx context.Context) {
	ch := b.lp.Run(ctx)
	for {
		select {
		case e := <-ch:
			b.handle(ctx, e)
		case <-b.lp.Done():
			b.close(b.lp.Err())
			return
		}
	}
}
