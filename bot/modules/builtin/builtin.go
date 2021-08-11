package builtin

import (
	"fmt"
	"time"

	"github.com/Toffee-iZt/HwBot/bot"
	"github.com/Toffee-iZt/HwBot/logger"
	"github.com/Toffee-iZt/HwBot/vkapi"
)

var Module = bot.Module{
	Name: "builtin",
	Init: func(_ *bot.Bot, l *logger.Logger) bool {
		start = time.Now()
		log = l
		return true
	},
	Terminate: nil,
	Commands:  []*bot.Command{&help, &about},
}

var log *logger.Logger

var start time.Time

var help = bot.Command{
	Cmd:  "help",
	Desc: "Помощь по командам",
	Help: "",
	Priv: true,
	Chat: true,
	Run: func(b *bot.Bot, msg *bot.IncomingMessage, a []string) {
		mods := b.Modules()
		if len(a) > 0 {
			for _, m := range mods {
				for _, c := range m.Commands {
					if c.Cmd == a[0] {
						b.SimpleReply(msg, fmt.Sprintf("%s - %s\n%s", c.Cmd, c.Desc, c.Help))
						return
					}
				}
			}
		}
		str := ">> Список модулей и команд\n"
		for _, m := range mods {
			str += "\n" + m.Name + "\n"
			for i, c := range m.Commands {
				if i > 0 {
					str += "\n"
				}
				str += "->" + c.Cmd + " - " + c.Desc
			}
		}
		b.SimpleReply(msg, str)
	},
}

var about = bot.Command{
	Cmd:  "about",
	Desc: "Информация о боте",
	Help: "",
	Priv: true,
	Chat: true,
	Run: func(b *bot.Bot, msg *bot.IncomingMessage, _ []string) {
		b.SimpleReply(msg, fmt.Sprintln(
			"HwBot\nVersion: debug\nSource code: github.com/Toffee-iZt/HwBot\n",
			"\n\nStart time:", start.Format("02 Jan 2006 15:04:05"),
			"\nUptime:", time.Now().Sub(start).Truncate(time.Second),
			"\nВерсия API VK:", vkapi.Version,
		))
	},
}
