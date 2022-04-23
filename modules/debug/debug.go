package debug

import (
	"embed"
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/Toffee-iZt/HwBot/bot"
	"github.com/Toffee-iZt/HwBot/bot/conversation"
	"github.com/Toffee-iZt/HwBot/common/format"
	"github.com/Toffee-iZt/HwBot/common/rt"
	"github.com/Toffee-iZt/HwBot/common/strbytes"
	"github.com/Toffee-iZt/HwBot/vkapi"
	"github.com/Toffee-iZt/HwBot/vkapi/keyboard"
)

// Module ...
var Module = bot.Module{
	Name:     "debug",
	Callback: onCallback,
	Commands: []*bot.Command{
		&debug,
	},
}

var debug = bot.Command{
	Cmd:         []string{"debug"},
	Description: "debug",
	Help:        "/debug [gc] - информация об используемой памяти (аргумент gc запустит полную сборку мусора)",
	Options:     bot.OptionInChat | bot.OptionInDialog,
	Run: func(ctx *bot.Context, msg *bot.Message) {
		var gc bool
		if len(msg.Args) > 0 {
			switch msg.Args[0] {
			case "obj":
				ctx.ReplyText(marshalObj(msg))
			case "embed":
				testembed(ctx)
			case "kb":
				testkb(ctx, msg.Args[1:])
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

func testembed(ctx *bot.Context) {
	f, err := resFs.Open("resources/gopher.png")
	if err != nil {
		ctx.ReplyText(err.Error())
	}

	ph, err := conversation.NewPhotoFile(f)
	if err != nil {
		ctx.ReplyText(err.Error())
	}
	f.Close()
	ctx.Reply(conversation.Message{
		Photos: []*conversation.Photo{ph},
	})
}

func testkb(ctx *bot.Context, a []string) {
	inline := strbytes.Has(a, "inline")
	onetime := strbytes.Has(a, "onetime")
	clear := strbytes.Has(a, "clear")

	if clear {
		kb := keyboard.Empty()
		ctx.Reply(conversation.Message{
			Text:     "keyboard cleared",
			Keyboard: kb,
		})
	}

	k := keyboard.New(onetime, inline)
	k.NextRow()
	k.AddText("LABEL", keyboard.ColorPositive, ctx.Payload("text"))
	k.AddCallback("CB", keyboard.ColorPrimary, ctx.Payload("cb"))
	k.AddCallback("Remove", keyboard.ColorNegative, ctx.Payload("edit"))

	ctx.Reply(conversation.Message{
		Text:     "testing keyborad",
		Keyboard: k,
	})
}

func onCallback(ctx *bot.Callback, payload vkapi.JSONData) {
	var a string
	err := payload.Unmarshal(&a)
	if err != nil {
		ctx.Log().Error("callback json ivalid")
		ctx.Close()
	}

	switch a {
	case "cb":
		ctx.ShowSnackbar("cool")
	case "edit":
		ctx.Edit(conversation.Message{
			Text: "good",
		})
	}
	ctx.Close()
}
