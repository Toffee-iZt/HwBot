package debug

import (
	"embed"
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/Toffee-iZt/HwBot/bot"
	"github.com/Toffee-iZt/HwBot/common/format"
	"github.com/Toffee-iZt/HwBot/common/rt"
	"github.com/Toffee-iZt/HwBot/vkapi"
)

// Module ...
var Module = bot.Module{
	Name:     "debug",
	Callback: onCallback,
	Commands: []*bot.Command{
		&debug,
		&testembed,
		&keyboard,
	},
}

var debug = bot.Command{
	Cmd:         []string{"debug"},
	Description: "debug",
	Help:        "/debug [gc] - информация об используемой памяти (аргумент gc запустит полную сборку мусора)",
	Options:     bot.OptionInChat | bot.OptionInDialog,
	Run: func(ctx *bot.Context, msg *bot.NewMessage, a []string) {
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
	Options:     bot.OptionInChat | bot.OptionInDialog,
	Run: func(ctx *bot.Context, msg *bot.NewMessage, a []string) {
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

var keyboard = bot.Command{
	Cmd:         []string{"kb"},
	Description: "test keyboard",
	Help:        "",
	Options:     bot.OptionInChat | bot.OptionInDialog,
	Run: func(ctx *bot.Context, msg *bot.NewMessage, a []string) {
		var inline, onetime bool
		if len(a) > 0 {
			inline = a[0] == "inline"
		}
		if len(a) > 1 {
			onetime = a[1] == "onetime"
		}

		kb := ctx.KeyboardGenerator(onetime, inline)
		kb.NextRow()
		kb.AddText("LABEL", vkapi.KeyboardColorPositive, `{"tap": "text"}`)
		kb.AddCallback("CB", vkapi.KeyboardColorPrimary, `{"tap": "cb"}`)

		ctx.Reply(bot.Message{
			Text:     "testing keyborad",
			Keyboard: kb,
		})
	},
}

func onCallback(ctx *bot.CallbackContext, payload vkapi.JSONData) {
	var a struct {
		Tap string `json:"tap"`
	}

	err := json.Unmarshal([]byte(payload), &a)
	if err != nil {
		ctx.Log().Error("callback json ivalid")
		ctx.ReplyCallback(nil)
	}

	if a.Tap == "" {
		ctx.ReplyCallback(nil)
	}

	ctx.ReplyCallback(&vkapi.EventData{
		Type: vkapi.EventDataTypeShowSnackbar,
		Text: "cool",
	})
}
