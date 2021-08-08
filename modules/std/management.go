package main

import (
	"HwBot/bot"
	"HwBot/core"
	"HwBot/vk"
	"HwBot/vkapi/vktypes"
)

var Module = bot.CommonModule{
	ModName: "Management",
	CmdList: nil,
	Init: func(b *bot.Bot) {
		// oninvite, onkick, onreturn, onleave
		/*
			bot.OnInvite = func(chat int, from, member int) {
				vk.SendText(chat, "Здарова").Assert("management OnInvite")
			}
			bot.OnReturn = func(chat int, member int) {
				vk.SendText(chat, "Ну и че ты вернулся").Assert("management OnReturn")
			}
			bot.OnKick = func(chat int, from, member int) {
				vk.SendText(chat, "Пока-пока").Assert("management OnKick")
			}
			bot.OnLeave = func(chat int, member int) {
				vk.SendText(chat, "Пока-пока").Assert("management OnLeave")
			}
		*/
	},
	Final: func() {},
}

var kick = bot.Command{
	Cmd:  "kick",
	Desc: "Исключить пользователя",
	Help: "...",
	Conv: bot.TypeChat,
	Run: func(b *bot.Bot, msg *vktypes.Message, args bot.Args) {
		vk := b.API()
		_ = vk.Kick(uint(msg.PeerID-2000000000), msg.Reply.FromID)
		//if err != nil {
		//	switch err.Code {
		//	case 935:
		//		err = vk.SendText(message.PeerID, "Данного пользователя нет в беседе")
		//	case 917:
		//		err = vk.SendText(message.PeerID, "Бот не является администратором")
		//	}
		//}

		return
	},
}

var by = bot.Command{
	Str:   "by",
	Short: "Кто пригласил пользователя",
	Help:  "...",
	Conv:  bot.TypeChat,
	Run: func(message *vk.Message, args bot.Args) (err *core.Error) {
		if message.Reply == nil {
			return
		}
		members, err := vk.GetChatMembers(uint(message.PeerID-2000000000), true)
		if err != nil {
			return vk.SendText(message.PeerID, "Бот не является администратором")
		}
		m := members.List[message.Reply.FromID]
		return vk.SendText(message.PeerID, vk.Mention(m.ID, "Пользователя")+" пригласил "+members.List[m.InvitedBy].Mention(false))
	},
}
