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
)

// Module ...
var Module = bot.Module{
	Name: "debug",
	Init: func(l *logger.Logger) bool {
		log = l
		log.Info("log test")
		return true
	},
	Terminate: nil,
	Commands: []*bot.Command{
		&debug,
		&testembed,
		&vkslow,
		//&keyboard,
	},
}

var log *logger.Logger

var debug = bot.Command{
	Cmd:         []string{"debug"},
	Description: "debug",
	Help:        "/debug [gc] - информация об используемой памяти (аргумент gc запустит полную сборку мусора)",
	InChat:      true,
	InPrivate:   true,
	Run: func(ctx *bot.MessageContext, msg *bot.NewMessage, a []string) {
		var gc bool
		if len(a) > 0 {
			switch a[0] {
			case "obj":
				ctx.ReplyText(marshalObj(msg))
			case "gc":
				gc = true
			}
		}
		ctx.ReplyText(memStats(gc))
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

//go:embed resources
var resFs embed.FS

var testembed = bot.Command{
	Cmd:         []string{"embed"},
	Description: "Test embed",
	Help:        "",
	InChat:      true,
	InPrivate:   true,
	Run: func(ctx *bot.MessageContext, msg *bot.NewMessage, a []string) {
		f, err := resFs.Open("resources/gopher.png")
		if err != nil {
			ctx.ReplyText(err.Error())
		}
		att, err := bot.NewAttachmentFile(f)
		f.Close()
		if err != nil {
			ctx.ReplyText(err.Error())
		}
		ctx.ReplyMessage("", att)
	},
}

var vkslow = bot.Command{
	Cmd:         []string{"vkslow"},
	Description: "",
	Help:        "",
	InChat:      true,
	InPrivate:   true,
	Run: func(ctx *bot.MessageContext, msg *bot.NewMessage, a []string) {
		n := time.Now().Unix()
		p := n - msg.Message.Date
		if p <= 2 {
			ctx.ReplyText("Ok. Vk time <2s")
			return
		}
		s := fmt.Sprintf(`Vk is slow.
		Time - %d
		Event time - %d
		Receive time - %d`, p, msg.Message.Date, n)
		ctx.ReplyText(s)
	},
}

/*
var keyboard = bot.Command{
	Cmd:  []string{"kb"},
	Description: "test keyboard",
	Help: "",
	InChat: true,
	InPrivate: true,
	Run: func(ctx *bot.Context, msg *bot.NewMessage, a []string) {
		var inline, onetime bool
		if len(a) > 0 {
			inline = a[0] == "inline"
		}
		if len(a) > 1 {
			onetime = a[1] == "onetime"
		}

		kb := vkapi.NewKeyboard(onetime, inline)
		kb.AddRow()
		kb.AddLocation(`{"tap": "location"}`)
		kb.AddRow()
		kb.AddText(`{"tap": "text"}`, "LABEL", vkapi.KeyboardColorPositive)
		kb.AddCallback(`{"tap": "cb"}`, "CB", vkapi.KeyboardColorPrimary)

		b.NewCallback(`{"tap": "cb"}`, func(c *bot.CallbackMessage) {
			b.API().Messages.SendMessageEventAnswer(
				c.EventID,
				c.UserID,
				c.PeerID,
				&vkapi.EventData{
					Type: vkapi.EventDataTypeShowSnackbar,
					Text: "cool",
				})
			log.Info("callback tapped")
		})

		_, err := b.API().Messages.Send(vkapi.OutMessage{
			PeerID:   m.Message.PeerID,
			Message:  "testing keyborad",
			Keyboard: kb,
		})
		if err != nil {
			log.Error(err.Error())
		}
	},
}
*/
