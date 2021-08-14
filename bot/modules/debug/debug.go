package debug

import (
	"embed"
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/Toffee-iZt/HwBot/bot"
	"github.com/Toffee-iZt/HwBot/common/format"
	"github.com/Toffee-iZt/HwBot/common/rt"
	"github.com/Toffee-iZt/HwBot/logger"
	"github.com/Toffee-iZt/HwBot/vkapi"
)

// Module ...
var Module = bot.Module{
	Name: "debug",
	Init: func(_ *bot.Bot, l *logger.Logger) bool {
		log = l
		log.Info("log test")
		return true
	},
	Terminate: nil,
	Commands: []*bot.Command{
		&debug,
		&ping,
		&testembed,
		&vkslow,
	},
}

var log *logger.Logger

var debug = bot.Command{
	Cmd:  "debug",
	Desc: "debug",
	Help: "/debug [gc] - информация об используемой памяти (аргумент gc запустит полную сборку мусора)",
	Chat: true,
	Priv: true,
	Run: func(b *bot.Bot, m *bot.IncomingMessage, args []string) {
		var gc bool
		if len(args) > 0 {
			switch args[0] {
			case "obj":
				b.SimpleReply(m, marshalObj(m))
				return
			case "gc":
				gc = true
			}
		}
		b.SimpleReply(m, memStats(gc))
	},
}

func marshalObj(obj interface{}) string {
	o, err := json.Marshal(obj)
	if err != nil {
		return err.Error()
	}
	return string(o)
}

func memStats(gc bool) string {
	ms := rt.GetMemStats(gc)

	var str string
	if gc {
		str = fmt.Sprintf("\n\nGarbage collection is done in %dns\nTotal stop-the-world time %dns", ms.GCTime, ms.PauseTotal)
	}
	str = fmt.Sprint(
		"Allocated:", format.FmtBytes(ms.Allocated, false, false),
		"\nHeap objects:", ms.Objects,
		"\nHeap inuse:", format.FmtBytes(ms.InUse, false, false),
		"\nGoroutines:", runtime.NumGoroutine(),
		str)
	return str
}

var ping = bot.Command{
	Cmd:  "ping",
	Desc: "Проверка работоспособности бота",
	Help: "",
	Chat: true,
	Priv: true,
	Run: func(b *bot.Bot, m *bot.IncomingMessage, _ []string) {
		b.SimpleReply(m, "понг")
	},
}

//go:embed resources
var resFs embed.FS

var testembed = bot.Command{
	Cmd:  "embed",
	Desc: "Test embed",
	Help: "",
	Chat: true,
	Priv: true,
	Run: func(b *bot.Bot, m *bot.IncomingMessage, _ []string) {
		f, err := resFs.Open("resources/gopher.png")
		if err != nil {
			b.SimpleReply(m, err.Error())
			return
		}
		vk := b.API()
		s, err := vk.UploadMessagesPhoto(m.Message.PeerID, f)
		if err != nil {
			log.Error("debug upload photo error: %s", err.Error())
			return
		}
		_, vkerr := vk.Messages.Send(vkapi.OutMessage{
			PeerID:     m.Message.PeerID,
			Attachment: []string{s},
		})
		if vkerr != nil {
			log.Error("debug send photo error: %s", vkerr.Error())
		}
	},
}

var vkslow = bot.Command{
	Cmd:  "vkslow",
	Desc: "",
	Help: "",
	Chat: true,
	Priv: true,
	Run: func(b *bot.Bot, m *bot.IncomingMessage, _ []string) {
		n := time.Now().Unix()
		p := n - m.Message.Date
		if p <= 2 {
			b.SimpleReply(m, "Ok. Vk time <2s")
			return
		}
		s := fmt.Sprintf(`Vk is slow.
		Time - %d
		Event time - %d
		Receive time - %d`, p, m.Message.Date, n)
		b.SimpleReply(m, s)
	},
}
