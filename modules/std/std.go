package std

import (
	"fmt"
	"strings"
	"time"

	"github.com/Toffee-iZt/HwBot/bot"
	"github.com/Toffee-iZt/HwBot/common/strbytes"
	"github.com/Toffee-iZt/HwBot/vkapi"
)

var (
	start    time.Time
	instance *bot.Bot
)

// Setup setups std module.
func Setup(b *bot.Bot) *bot.Module {
	start = time.Now()
	instance = b
	return &module
}

var module = bot.Module{
	Name: "std",
	Commands: []*bot.Command{
		&help,
		&about,
		&stat,
	},
}

var help = bot.Command{
	Cmd:         []string{"help", "помощь", "хелп"},
	Description: "Помощь по командам",
	Help:        "",
	Options:     bot.OptionInChat | bot.OptionInDialog,
	Run: func(ctx *bot.Context, msg *bot.Message) {
		mods := instance.ModList()
		if len(msg.Args) > 0 {
			for _, m := range mods {
				for _, c := range m.Commands {
					if strbytes.Has(c.Cmd, msg.Args[0]) {
						ctx.ReplyText(fmt.Sprintf("%s - %s\nAliases: %s\n\n%s", c.Cmd[0], c.Description, strings.Join(c.Cmd, ", "), c.Help))
					}
				}
			}
		}
		str := ">> Список модулей и команд\n"
		for _, m := range mods {
			str += "\n"
			for i, c := range m.Commands {
				if i > 0 {
					str += "\n"
				}
				str += " -> " + c.Cmd[0] + " - " + c.Description
			}
		}
		ctx.ReplyText(str)
	},
}

var about = bot.Command{
	Cmd:         []string{"about"},
	Description: "Информация о боте",
	Help:        "",
	Options:     bot.OptionInChat | bot.OptionInDialog,
	Run: func(ctx *bot.Context, msg *bot.Message) {
		ctx.ReplyText(fmt.Sprintln(
			"HwBot\nVersion: debug\nSource code: github.com/Toffee-iZt/HwBot\n",
			"\nВерсия API VK:", vkapi.Version,
		))
	},
}

var stat = bot.Command{
	Cmd:         []string{"stat", "ping"},
	Description: "Статистика",
	Help:        "Проверка работоспособности бота и статистика работы",
	Options:     bot.OptionInChat | bot.OptionInDialog,
	Run: func(ctx *bot.Context, msg *bot.Message) {
		p := time.Now().Unix() - msg.Message.Date
		str := fmt.Sprintf(`Stats
		
		Start time: %s
		Uptime: %s`, start.Format("02 Jan 2006 15:04:05"), time.Now().Sub(start).Truncate(time.Second))

		if p > 3 {
			str += fmt.Sprint("\n\nVK event delay:", p)
		}

		ctx.ReplyText(str)
	},
}
