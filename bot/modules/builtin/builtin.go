package builtin

import (
	"fmt"
	"strings"
	"time"

	"github.com/Toffee-iZt/HwBot/bot"
	"github.com/Toffee-iZt/HwBot/common/strbytes"
	"github.com/Toffee-iZt/HwBot/logger"
	"github.com/Toffee-iZt/HwBot/vkapi"
)

var Module = bot.Module{
	Name: "builtin",
	Init: func(l *logger.Logger) bool {
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
	Cmd:         []string{"help", "помощь", "хелп"},
	Description: "Помощь по командам",
	Help:        "",
	InPrivate:   true,
	InChat:      true,
	Run: func(ctx *bot.Context, msg *bot.NewMessage, a []string) {
		mods := ctx.BotInstance().ModList()
		if len(a) > 0 {
			for _, m := range mods {
				for _, c := range m.Commands {
					if strbytes.Has(c.Cmd, a[0]) {
						ctx.ReplyText(fmt.Sprintf("%s - %s\nAliases: %s\n\n%s", c.Cmd[0], c.Description, strings.Join(c.Cmd, ", "), c.Help))
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
				str += "->" + c.Cmd[0] + " - " + c.Description
			}
		}
		ctx.ReplyText(str)
	},
}

var about = bot.Command{
	Cmd:         []string{"about"},
	Description: "Информация о боте",
	Help:        "",
	InPrivate:   true,
	InChat:      true,
	Run: func(ctx *bot.Context, msg *bot.NewMessage, _ []string) {
		ctx.ReplyText(fmt.Sprintln(
			"HwBot\nVersion: debug\nSource code: github.com/Toffee-iZt/HwBot\n",
			"\n\nStart time:", start.Format("02 Jan 2006 15:04:05"),
			"\nUptime:", time.Now().Sub(start).Truncate(time.Second),
			"\nВерсия API VK:", vkapi.Version,
		))
	},
}
