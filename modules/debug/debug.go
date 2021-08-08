package main

import (
	botapi "HwBot/bot/api"
	"HwBot/common/format"
	"HwBot/common/rt"
	"HwBot/vkapi"
	"embed"
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

// Module ...
var Module = botapi.Module{
	Name: "Debug",
	Init: func(_ botapi.Bot, l botapi.Logger) bool {
		log = l
		log.Info("debug successfully inited")
		return true
	},
	Terminate: nil,
	Commands: []*botapi.Command{
		&debug,
		&ping,
		&testembed,
		&vkslow,
	},
	OnMessage: func(b botapi.Bot, msg *botapi.IncomingMessage) {
		m := &msg.Message
		u, err := b.API().Users.Get(m.FromID)
		if err != nil {
			log.Error(err.Error())
			return
		}
		log.Info(`New message:
		from: %s %s
		message: %s`, u[0].FirstName, u[0].LastName, m.Text)
	},
}

var log botapi.Logger

var debug = botapi.Command{
	Cmd:  "debug",
	Desc: "debug",
	Help: "/debug [gc] - информация об используемой памяти (аргумент gc запустит полную сборку мусора)",
	Chat: true,
	Priv: true,
	Run: func(b botapi.Bot, m *botapi.IncomingMessage, args []string) {
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

var ping = botapi.Command{
	Cmd:  "ping",
	Desc: "Проверка работоспособности бота",
	Help: "",
	Chat: true,
	Priv: true,
	Run: func(b botapi.Bot, m *botapi.IncomingMessage, _ []string) {
		b.SimpleReply(m, "понг")
	},
}

//go:embed resources
var resFs embed.FS

var testembed = botapi.Command{
	Cmd:  "embed",
	Desc: "Test embed",
	Help: "",
	Chat: true,
	Priv: true,
	Run: func(b botapi.Bot, m *botapi.IncomingMessage, _ []string) {
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
		_, vkerr := vk.Messages.Send(
			vkapi.MessagePeer{
				PeerID: m.Message.PeerID,
			},
			vkapi.MessageContent{
				Attachment: []string{s},
			},
		)
		if vkerr != nil {
			log.Error("debug send photo error: %s", vkerr.Error())
		}
	},
}

var vkslow = botapi.Command{
	Cmd:  "vkslow",
	Desc: "",
	Help: "",
	Chat: true,
	Priv: true,
	Run: func(b botapi.Bot, m *botapi.IncomingMessage, _ []string) {
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
