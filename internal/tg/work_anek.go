package tg

import (
	"log"
	"strings"

	"github.com/bochkov/gobot/internal/anekdot"
)

type AnekdotWorker struct {
	anekdot.Service
}

func NewAnekdotWorker(s anekdot.Service) *AnekdotWorker {
	return &AnekdotWorker{Service: s}
}

func (a *AnekdotWorker) Description() string {
	return a.Service.Description()
}

func (a *AnekdotWorker) IsMatch(text string) bool {
	return strings.Contains(strings.ToLower(text), "анек") || strings.Contains(strings.ToLower(text), "joke")
}

func (a *AnekdotWorker) Answer(chatId int64, txt string) []Method {
	sm := SendMessage[int64]{ChatId: chatId}
	anek, err := a.Service.GetRandom()
	if err != nil {
		log.Print(err)
		sm.Text = "не смог найти анекдот"
	} else {
		sm.Text = anek.Text
	}
	res := make([]Method, 0)
	res = append(res, &sm)
	return res
}
