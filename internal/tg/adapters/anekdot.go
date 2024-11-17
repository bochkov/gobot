package adapters

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bochkov/gobot/internal/services/anekdot"
	"github.com/bochkov/gobot/internal/tg"
)

type AnekdotAdapter struct {
	service anekdot.Service
}

func NewAnekdotAdapter(s anekdot.Service) AnekdotAdapter {
	return AnekdotAdapter{service: s}
}

func (a AnekdotAdapter) Description() string {
	return a.service.Description()
}

func (a AnekdotAdapter) IsMatch(text string) bool {
	return strings.Contains(strings.ToLower(text), "анек") || strings.Contains(strings.ToLower(text), "joke")
}

func (a AnekdotAdapter) Answer(chatId int64, txt string) []tg.Method {
	receivers := []string{strconv.FormatInt(chatId, 10)}
	methods, err := a.PushData(receivers)
	if err != nil {
		sm := tg.SendMessage[int64]{ChatId: chatId}
		sm.Text = "не получилось ("
		return []tg.Method{&sm}
	}
	return methods
}

func (a AnekdotAdapter) PushData(receivers []string) ([]tg.Method, error) {
	anek, err := a.service.GetRandom()
	if err != nil {
		return nil, err
	}
	res := make([]tg.Method, 0)
	for _, recv := range receivers {
		sm := new(tg.SendMessage[string])
		sm.ChatId = recv
		sm.Text = fmt.Sprintf("Анекдот дня:\n%s", anek.Text)
		res = append(res, sm)
	}
	return res, nil
}
