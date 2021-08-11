package random

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/Toffee-iZt/HwBot/bot"
	"github.com/Toffee-iZt/HwBot/common/strbytes"
	"github.com/Toffee-iZt/HwBot/logger"
	"github.com/Toffee-iZt/HwBot/vkapi/vkutils"
)

// Module ...
var Module = bot.Module{
	Name: "random",
	Init: func(_ *bot.Bot, _ *logger.Logger) bool {
		myRand = rand.New(rand.NewSource(time.Now().UnixNano()))
		return true
	},
	Commands: []*bot.Command{&list, &number, &who, &flip, &info, &when},
}

var myRand *rand.Rand

var list = bot.Command{
	Cmd:  "list",
	Desc: "Рандомный список участников",
	Help: "/список умных - список из 5 рандомных участников" +
		"\nПосле команды можно указать длину списка" +
		"\nНапример /список 12 задротов - список из 12 участников",
	Chat: true,
	Priv: false,
	Run: func(b *bot.Bot, m *bot.IncomingMessage, args []string) {
		members, err := b.API().Messages.GetChatMembers(vkutils.PeerToChat(m.Message.PeerID))
		if err != nil {
			b.SimpleReply(m, "У бота недостаточно прав доступа для выполнения команды")
			return
		}

		l := len(members.Profiles)

		var num = 5
		if len(args) > 0 {
			n, _ := strconv.Atoi(args[0])
			if n > 0 {
				num = n
			}
		}
		if num > l {
			num = l
		}

		users := members.Profiles
		var str = "Список " + strbytes.Join(args, ' ') + ":\n"
		for i := 0; i < num; i++ {
			n := myRand.Intn(l)
			u := &users[n]
			str += strconv.Itoa(i+1) + ". " + vkutils.Mention(u.ID, u.FirstName+" "+u.LastName) + "\n"
			l--
			users[n] = users[l]
		}
		b.SimpleReply(m, str)
	},
}

var number = bot.Command{
	Cmd:  "rand",
	Desc: "Рандомное число",
	Help: "/rand - рандомное число [0, 100]" +
		"\n/rand <max> - рандомное число [0, max] (0 < max <= maxInt64)" +
		"\nmax также принимает 2(0b), 8(0 или 0o), 16(0x) системы счисления",
	Chat: true,
	Priv: true,
	Run: func(b *bot.Bot, m *bot.IncomingMessage, args []string) {
		var max int64 = 100
		var base = 10
		if len(args) > 0 {
			num := args[0]
			if num[0] == '0' {
				switch num[1] {
				case 'b':
					base = 2
					num = num[2:]
				case 'o':
					base = 8
					num = num[2:]
				case 'x':
					base = 16
					num = num[2:]
				default:
					base = 8
					num = num[1:]
				}
			}
			var err error
			max, err = strconv.ParseInt(num, base, 64)
			if err != nil || max <= 0 {
				b.SimpleReply(m, "Укажите число от 0 до MaxInt64")
				return
			}
		}
		r := myRand.Int63n(max)
		t := "Рандомное число — " + strconv.FormatInt(r, base)
		if base != 10 {
			t += "\n10 -> " + strconv.FormatInt(r, 10)
		}

		b.SimpleReply(m, t)
	},
}

var who = bot.Command{
	Cmd:  "кто",
	Desc: "Выбрать рандомного участника",
	Help: "/кто <string>",
	Chat: true,
	Priv: false,
	Run: func(b *bot.Bot, m *bot.IncomingMessage, args []string) {
		members, err := b.API().Messages.GetChatMembers(vkutils.PeerToChat(m.Message.PeerID))
		if err != nil {
			b.SimpleReply(m, "У бота недостаточно прав доступа для выполнения команды")
			return
		}

		r := myRand.Intn(len(members.Profiles))
		p := members.Profiles[r]
		str := vkutils.Mention(p.ID, p.FirstName+" "+p.LastName)
		if len(args) > 0 {
			str += " — " + strbytes.Join(args, ' ')
		}

		b.SimpleReply(m, str)
	},
}

var flip = bot.Command{
	Cmd:  "flip",
	Desc: "Подбросить монетку",
	Help: "",
	Chat: true,
	Priv: true,
	Run: func(b *bot.Bot, m *bot.IncomingMessage, _ []string) {
		var r string
		if myRand.Intn(2) == 1 {
			r = "Выпал орёл"
		} else {
			r = "Выпала решка"
		}
		b.SimpleReply(m, r)
	},
}

var info = bot.Command{
	Cmd:  "инфа",
	Desc: "Вероятность события",
	Help: "/инфа <событие> - случайная вероятность события",
	Chat: true,
	Priv: true,
	Run: func(b *bot.Bot, m *bot.IncomingMessage, args []string) {
		if len(args) == 0 {
			b.SimpleReply(m, "Укажите событие")
			return
		}
		p := myRand.Intn(101)
		e := strbytes.Join(args, ' ')
		b.SimpleReply(m, "Вероятность того, что "+e+" — "+strconv.Itoa(p)+"%")
	},
}

var when = bot.Command{
	Cmd:  "когда",
	Desc: "Когда произойдет событие",
	Help: "/когда <событие> - случайная дата события",
	Chat: true,
	Priv: true,
	Run: func(b *bot.Bot, m *bot.IncomingMessage, args []string) {
		if len(args) == 0 {
			b.SimpleReply(m, "Укажите событие")
			return
		}
		t := time.Now().AddDate(myRand.Intn(51), myRand.Intn(12), myRand.Intn(31))
		e := strbytes.Join(args, ' ')
		b.SimpleReply(m, e+" "+t.Format("02 Jan 2006"))
	},
}
