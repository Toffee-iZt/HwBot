package bot

import (
	"strings"

	"github.com/Toffee-iZt/HwBot/vkapi/longpoll"
)

// NewMessage allias.
type NewMessage longpoll.MessageNew

// CommandPrefixes are the characters with which commands must begin.
const CommandPrefixes = "/!"

// MessageLimit is a limit to the number of characters of the input message.
const MessageLimit = 300

func (b *Bot) getCommand(msg *longpoll.Message) (*Command, []string) {
	if n := len(msg.Text); n > 1 && strings.IndexByte(CommandPrefixes, msg.Text[0]) != -1 {
		if n > MessageLimit {
			msg.Text = msg.Text[:MessageLimit]
		}
		s := strings.Split(msg.Text, " ")
		return b.commands[s[0][1:]], s[1:]
	}
	return nil, nil
}

func (b *Bot) onMessage(msg *longpoll.MessageNew) {
	cmd, args := b.getCommand(&msg.Message)
	if cmd == nil {
		return
	}
	bctx := &MessageContext{
		eventctx: eventctx{
			Conv: &Conversation{
				peer: msg.Message.PeerID,
				api:  b.vk,
			},
			bot: b,
			mod: cmd.module,
		},
	}

	id := int(msg.Message.PeerID)
	if id > 2e9 && !cmd.InChat || id < 2e9 && !cmd.InPrivate {
		bctx.ReplyText("Команда недоступна в данном типе чата")
	}

	cmd.Run(bctx, (*NewMessage)(msg), args)
}
