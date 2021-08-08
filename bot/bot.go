package bot

import (
	"HwBot/bot/api"
	"HwBot/common"
	"HwBot/common/execdir"
	"HwBot/logger"
	"HwBot/vkapi"
	"HwBot/vkapi/longpoll"
	"context"
	"os"
	"path/filepath"
	"plugin"
)

// Config is bot config.
type Config struct {
	Prefixes    string
	ModulesPath string
}

// New creates new bot.
func New(vk *vkapi.Client, cfg Config, w *logger.Writer) *Bot {
	return &Bot{
		log:      logger.New(w, "BOT"),
		cfg:      cfg,
		vk:       vk,
		lp:       longpoll.New(vk, 25),
		commands: make(map[string]*api.Command),
	}
}

// Bot is vk bot.
type Bot struct {
	sync common.Sync

	cfg Config

	log *logger.Logger
	vk  *vkapi.Client
	lp  *longpoll.LongPoll

	commands map[string]*api.Command
	modules  []*api.Module
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
func (b *Bot) Run(ctx context.Context) bool {
	if !b.sync.Init() {
		return false
	}

	b.log.Info("loading modules")

	n, err := b.loadModules()
	if err != nil {
		b.log.Error("modules folder reading error: %s", err.Error())
	}
	if n == 0 {
		b.log.Warn("0 modules have been loaded")
	} else {
		b.log.Info("%d modules have been loaded", n)
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

func (b *Bot) loadModules() (int, error) {
	const moduleFileExt = ".bmod"

	edir := execdir.GetExecDir()
	entries, err := execdir.ReadDir(b.cfg.ModulesPath)
	if os.IsNotExist(err) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	modLog := b.log.Child("MODULE")

	for _, i := range entries {
		name := i.Name()
		if i.IsDir() || filepath.Ext(name) != moduleFileExt {
			b.log.Info("skip no module: %s", name)
			continue
		}

		b.log.Info("loading module %s", name)

		p, err := plugin.Open(filepath.Join(edir, b.cfg.ModulesPath, name))
		if err != nil {
			b.log.Error(err.Error()) // already wrapper
			continue
		}

		modSymbol, _ := p.Lookup(api.ModuleSymbol)
		if modSymbol == nil {
			b.log.Error(`there is no "%s" symbol in module "%s"`, modSymbol, name)
			continue
		}

		mod, ok := modSymbol.(*api.Module)
		if !ok {
			b.log.Error(`"%s" symbol in module "%s" has wrong type`, modSymbol, name)
			continue
		}

		b.initModule(mod, modLog)
	}

	return len(b.modules), nil
}
