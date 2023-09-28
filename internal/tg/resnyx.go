package tg

import (
	"log"
	"strings"

	"github.com/bochkov/gobot/db"
)

type MethodCustomize func(message *SendMessage[string])

type PushService struct {
}

func (p *PushService) Push(text string, mc ...MethodCustomize) {
	if text == "" {
		log.Print("empty text")
		return
	}
	token := db.GetProp(db.TgBotTokenKey, "")
	if token == "" {
		log.Print("no token specified")
		return
	}
	chatId := db.GetProp(db.ChatAutoSend, "") // TODO
	if chatId == "" {
		log.Print("no chat.id specified")
		return
	}
	bot := NewBot(token)
	for _, chat := range strings.Split(chatId, ";") {
		sm := new(SendMessage[string])
		sm.ChatId = chat
		sm.Text = text
		for _, customize := range mc {
			customize(sm)
		}
		if _, exec := bot.Execute(sm); exec != nil {
			log.Print(exec)
		}
	}
}
