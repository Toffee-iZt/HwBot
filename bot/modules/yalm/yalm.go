package yalm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Toffee-iZt/HwBot/bot"
	"github.com/Toffee-iZt/HwBot/logger"
	"github.com/Toffee-iZt/balaboba"
)

// Module ...
var Module = bot.Module{
	Name: "yalm",
	Init: func(_ *bot.Bot, l *logger.Logger) bool {
		log = l
		client = balaboba.New()
		err := client.Options()
		if err != nil {
			log.Error("init error: %s", err.Error())
			return false
		}
		return true
	},
	Commands: []*bot.Command{
		&yalm,
	},
}

/////////////////////////////////////////////

var client *balaboba.Client

var log *logger.Logger

const pErr = "Произошла ошибка в работе бота. Разработчик уже знает об этом, в скором времени будет пофикшено."
const apiErrorFmt = "Произошла ошибка %s (%d) при работе с API балабобы. Такое случается, так как не изучены все аспекты работы API. Если вы знаете причину и как это исправить - напишите разработчику."

const warns = balaboba.Warn1 + "\n" + balaboba.Warn2
const help = `
/yalm warns - предупреждения (обязательно прочитайте перед использованием!)
/yalm {style} <string> - сгенерировать текст
стиль указывается в формате sN, где N - номер стиля (см. styles)
/yalm styles - список стилей генерации текста`

var yalm = bot.Command{
	Cmd:  "yalm",
	Desc: "yandex balaboba (yalm)",
	Help: balaboba.About + help,
	Chat: true,
	Priv: true,
	Run: func(b *bot.Bot, m *bot.IncomingMessage, a []string) {
		if len(a) == 0 {
			b.SimpleReply(m, help)
			return
		}
		switch a[0] {
		case "warns":
			b.SimpleReply(m, warns)
			return
		case "styles":
			intros, err := client.Intros()
			if err != nil {
				log.Error("yalm intros error: ", err.Error())
				b.SimpleReply(m, pErr)
				return
			}
			if intros.Error != 0 {
				log.Error("yalm intros response error: ", err.Error())
				b.SimpleReply(m, fmt.Sprintf(apiErrorFmt, "api/yalm/intros", intros.Error))
				return
			}
			msg := "Стили генерации текста\n"
			for i := range intros.Intros {
				s := intros.Intros[i]
				msg += fmt.Sprintf("\n%d - %s - %s", s.Style, s.String, s.Description)
			}
			b.SimpleReply(m, msg)
			return
		}

		var s balaboba.Style
		if a[0][0] == 's' {
			u, _ := strconv.ParseUint(a[0][1:], 10, 8)
			s = balaboba.Style(u)
			if s.Invalid() {
				s = balaboba.NoStyle
			}
			a = a[1:]
		}

		query := strings.Join(a, " ")
		if query == "" {
			b.SimpleReply(m, help)
			return
		}

		g, err := client.Get(query, s)
		if err != nil {
			log.Error("balaboba get error: %s", err.Error())
			b.SimpleReply(m, pErr)
			return
		}
		if g.Error != 0 {
			log.Error("balaboba get response error:", g.Error)
			b.SimpleReply(m, fmt.Sprintf(apiErrorFmt, "api/yalm/text3", g.Error))
			return
		}
		if g.BadQuery != 0 {
			b.SimpleReply(m, balaboba.BadQuery)
			return
		}
		b.SimpleReply(m, g.Query+" "+g.Text)
	},
}
