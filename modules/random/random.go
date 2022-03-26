package random

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/Toffee-iZt/HwBot/bot"
	"github.com/Toffee-iZt/HwBot/common/strbytes"
	"github.com/Toffee-iZt/HwBot/vkapi/vkutils"
)

// Module ...
var Module = bot.Module{
	Name: "random",
	Init: func() error {
		myRand = rand.New(rand.NewSource(time.Now().UnixNano()))
		return nil
	},
	Commands: []*bot.Command{&list, &number, &who, &flip, &info, &when},
}

var myRand *rand.Rand

var list = bot.Command{
	Cmd:         []string{"list", "список"},
	Description: "Рандомный список участников",
	Help: "/список умных - список из 5 рандомных участников" +
		"\nПосле команды можно указать длину списка" +
		"\nНапример /список 12 задротов - список из 12 участников",
	Options: bot.OptionInChat | bot.OptionInDialog,
	Run: func(ctx *bot.Context, msg *bot.NewMessage, a []string) {
		members, err := ctx.Conv.GetMembers()
		if err != nil {
			ctx.ReplyError("У бота недостаточно прав доступа для выполнения команды")
		}

		l := len(members.Profiles)

		var num = 5
		if len(a) > 0 {
			n, _ := strconv.Atoi(a[0])
			if n > 0 {
				num = n
			}
		}
		if num > l {
			num = l
		}

		users := members.Profiles
		var str = "Список " + strbytes.Join(a, ' ') + ":\n"
		for i := 0; i < num; i++ {
			n := myRand.Intn(l)
			u := users[n]
			str += strconv.Itoa(i+1) + ". " + vkutils.Mention(u.ID.ToID(), u.FirstName+" "+u.LastName) + "\n"
			l--
			users[n] = users[l]
		}
		ctx.ReplyText(str)
	},
}

var number = bot.Command{
	Cmd:         []string{"rand", "roll"},
	Description: "Рандомное число",
	Help: "/rand - рандомное число [0, 100]" +
		"\n/rand <max> - рандомное число [0, max] (0 < max <= maxInt64)",
	Options: bot.OptionInChat | bot.OptionInDialog,
	Run: func(ctx *bot.Context, msg *bot.NewMessage, a []string) {
		var max int64 = 100
		if len(a) > 0 {
			num := a[0]
			var err error
			max, err = strconv.ParseInt(num, 10, 64)
			if err != nil || max <= 0 {
				ctx.ReplyText("Укажите число от 0 до MaxInt64")
			}
		}
		ctx.ReplyText("Выпало число " + strconv.FormatInt(myRand.Int63n(max), 10))
	},
}

var who = bot.Command{
	Cmd:         []string{"who", "кто"},
	Description: "Выбрать рандомного участника",
	Help:        "/who <string>",
	Options:     bot.OptionInChat | bot.OptionInDialog,
	Run: func(ctx *bot.Context, msg *bot.NewMessage, a []string) {
		members, err := ctx.Conv.GetMembers()
		if err != nil {
			ctx.ReplyText("У бота недостаточно прав доступа для выполнения команды")
			return
		}

		r := myRand.Intn(len(members.Profiles))
		p := members.Profiles[r]
		str := vkutils.Mention(p.ID.ToID(), p.FirstName+" "+p.LastName)
		if len(a) > 0 {
			str += " — " + strbytes.Join(a, ' ')
		}

		ctx.ReplyText(str)
	},
}

var flip = bot.Command{
	Cmd:         []string{"flip", "флип", "монетка"},
	Description: "Подбросить монетку",
	Help:        "",
	Options:     bot.OptionInChat | bot.OptionInDialog,
	Run: func(ctx *bot.Context, msg *bot.NewMessage, a []string) {
		var r string
		if myRand.Intn(2) == 1 {
			r = "Выпал орёл"
		} else {
			r = "Выпала решка"
		}
		ctx.ReplyText(r)
	},
}

var info = bot.Command{
	Cmd:         []string{"info", "инфа"},
	Description: "Вероятность события",
	Help:        "/info <событие> - случайная вероятность события",
	Options:     bot.OptionInChat | bot.OptionInDialog,
	Run: func(ctx *bot.Context, msg *bot.NewMessage, a []string) {
		if len(a) == 0 {
			ctx.ReplyText("Укажите событие")
			return
		}
		p := myRand.Intn(101)
		e := strbytes.Join(a, ' ')
		ctx.ReplyText("Вероятность того, что " + e + " — " + strconv.Itoa(p) + "%")
	},
}

var when = bot.Command{
	Cmd:         []string{"when", "когда"},
	Description: "Когда произойдет событие",
	Help:        "/when <событие> - случайная дата события",
	Options:     bot.OptionInChat | bot.OptionInDialog,
	Run: func(ctx *bot.Context, msg *bot.NewMessage, a []string) {
		if len(a) == 0 {
			ctx.ReplyText("Укажите событие")
			return
		}
		t := time.Now().AddDate(myRand.Intn(51), myRand.Intn(12), myRand.Intn(31))
		e := strbytes.Join(a, ' ')
		ctx.ReplyText(e + " " + t.Format("02 Jan 2006"))
	},
}
