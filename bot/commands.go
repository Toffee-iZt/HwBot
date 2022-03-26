package bot

import (
	"strings"

	"github.com/Toffee-iZt/HwBot/common"
	"github.com/Toffee-iZt/HwBot/vkapi/longpoll"
)

// NewMessage allias.
type NewMessage longpoll.MessageNew

// CommandPrefixes are the characters with which commands must begin.
const CommandPrefixes = "/!"

// Command respresents conversation command.
type Command struct {
	Run         func(*Context, *NewMessage, []string)
	Cmd         []string
	Description string
	Help        string
	Options     common.Flag
	module      *Module
}

// Command options.
const (
	OptionInDialog = 1 << iota
	OptionInChat
)

func (b *Bot) onMessage(msg *longpoll.MessageNew) {
	// MessageLimit is a limit to the number of characters of the input message.
	const MessageLimit = 300

	n := len(msg.Message.Text)
	if n == 0 || strings.IndexByte(CommandPrefixes, msg.Message.Text[0]) == -1 {
		return
	}
	if n > MessageLimit {
		msg.Message.Text = msg.Message.Text[:MessageLimit]
	}

	s := strings.Split(msg.Message.Text, " ")
	cmd, args := b.commands[s[0][1:]], s[1:]
	if cmd == nil {
		return
	}

	bctx := &Context{
		eventctx: eventctx{
			Conv: &Conversation{
				peer: msg.Message.PeerID,
				api:  b.vk,
			},
			mod: cmd.module,
		},
	}

	id := int(msg.Message.PeerID)
	if id > 2e9 && !cmd.Options.Has(OptionInChat) || id < 2e9 && !cmd.Options.Has(OptionInDialog) {
		bctx.ReplyText("Команда недоступна в данном типе чата")
	}

	cmd.Run(bctx, (*NewMessage)(msg), args)
}
