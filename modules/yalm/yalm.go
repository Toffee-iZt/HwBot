package yalm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Toffee-iZt/HwBot/bot"
	"github.com/Toffee-iZt/balaboba"
)

// Module ...
var Module = bot.Module{
	Name: "yalm",
	Commands: []*bot.Command{
		&yalm,
	},
}

/////////////////////////////////////////////

var client = balaboba.New()
var isAvailable = client.IsAvailable()

const pErr = "Произошла ошибка в работе бота. Разработчик уже знает об этом, в скором времени будет пофикшено."
const apiErrorFmt = "Произошла ошибка %s (%d) при работе с API балабобы. Такое случается, так как не изучены все аспекты работы API. Если вы знаете причину и как это исправить - напишите разработчику."

const warns = balaboba.Warn1 + "\n" + balaboba.Warn2
const help = `
/yalm warns - предупреждения (обязательно прочитайте перед использованием!)
/yalm {style} <string> - сгенерировать текст
стиль указывается в формате sN, где N - номер стиля (см. styles)
/yalm styles - список стилей генерации текста`

var yalm = bot.Command{
	Cmd:         []string{"yalm"},
	Description: "yandex balaboba (yalm)",
	Help:        balaboba.About + help,
	Options:     bot.OptionInChat | bot.OptionInDialog,
	Run: func(ctx *bot.Context, msg *bot.Message) {
		if !isAvailable {
			ctx.ReplyText("Сервис временно недоступен")
		}
		if len(msg.Args) == 0 {
			ctx.ReplyText(help)
		}
		switch msg.Args[0] {
		case "warns":
			ctx.ReplyText(warns)
		case "styles":
			styles(ctx)
		}

		var s balaboba.Style
		if msg.Args[0][0] == 's' {
			u, _ := strconv.ParseUint(msg.Args[0][1:], 10, 8)
			s = balaboba.Style(u)
			msg.Args = msg.Args[1:]
		}

		query := strings.Join(msg.Args, " ")
		if query == "" {
			ctx.ReplyText(help)
		}

		generate(ctx, query, s)
	},
}

func styles(ctx *bot.Context) {
	log := ctx.Log()
	intros, err := client.Intros()
	if err != nil {
		log.Error("yalm intros error: ", err.Error())
		ctx.ReplyText(pErr)
	}
	if intros.Error != 0 {
		log.Error("yalm intros response error: ", err.Error())
		ctx.ReplyText(fmt.Sprintf(apiErrorFmt, "api/yalm/intros", intros.Error))
	}
	msg := "Стили генерации текста\n"
	for i := range intros.Intros {
		s := intros.Intros[i]
		msg += fmt.Sprintf("\n%d - %s - %s", s.Style, s.String, s.Description)
	}
	ctx.ReplyText(msg)
}

func generate(ctx *bot.Context, query string, s balaboba.Style) {
	log := ctx.Log()
	g, err := client.Get(nil, query, s)
	if err != nil {
		log.Error("balaboba get error: %s", err.Error())
		ctx.ReplyText(pErr)
	}
	if g.Error != 0 {
		log.Error("balaboba get response error:", g.Error)
		ctx.ReplyText(fmt.Sprintf(apiErrorFmt, "api/yalm/text3", g.Error))
	}
	if g.BadQuery != 0 {
		ctx.ReplyText(balaboba.BadQuery)
	}
	ctx.ReplyText(g.Query + " " + g.Text)
}
