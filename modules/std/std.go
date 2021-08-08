package main

import (
	"HwBot/bot"
	"HwBot/core"
	"HwBot/vk"
	"fmt"
	"time"
)

// Module ...
var Module = bot.Module{
	Name:  "Std",
	Init:  nil,
	Final: nil,
	List: []*bot.Command{
		&about,
		&help,
	},
}

var about = bot.Command{
	Str:   "about",
	Short: "О боте",
	Help:  "Информация о боте",
	Conv:  bot.TypeChat | bot.TypePrivate,
	Run: func(message *vk.Message, args bot.Args) *core.Error {
		return vk.SendText(message.PeerID, fmt.Sprintln(
			"HwBot\nVersion: core tests\n",
			"Написан на GoLang с библиотеками fastjson и fasthttp",
			"\n\nStart time:", bot.StartTime.Format("02 Jan 2006 15:04:05"),
			"\nUptime:", time.Now().Sub(bot.StartTime).Truncate(time.Second),
			"\nВерсия API VK:", vk.APIVersion(),
		))
	},
}

var help = bot.Command{
	Str:   "help",
	Short: "Помощь по модулям и командам",
	Help:  "",
	Conv:  bot.TypeChat | bot.TypePrivate,
	Run: func(message *vk.Message, args bot.Args) *core.Error {
		if len(args) > 0 {
			for _, m := range bot.Modules {
				if name := m.Name; name == args[0] {
					str := "Модуль " + name + "\n"
					for _, c := range m.List {
						str += "\n" + c.Str + " - " + c.Short
					}
					return vk.SendText(message.PeerID, str)
				}
			}
		}
		str := ">> Список модулей и команд\n"
		for _, m := range bot.Modules {
			str += "\n" + m.Name + "\n—>"
			for i, c := range m.List {
				if i > 0 {
					str += ", "
				}
				str += c.Str
			}
		}
		return vk.SendText(message.PeerID, str)
	},
}
