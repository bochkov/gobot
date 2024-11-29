package adapters

import (
	"strconv"
	"strings"

	"github.com/bochkov/gobot/internal/services/wiki"
	"github.com/bochkov/gobot/internal/tg"
)

type WPotdAdapter struct {
	service *wiki.PotdService
}

func NewWPotdAdapter(s *wiki.PotdService) *WPotdAdapter {
	return &WPotdAdapter{service: s}
}

func (w *WPotdAdapter) Description() string {
	return w.service.Description()
}

func (w *WPotdAdapter) IsMatch(text string) bool {
	return strings.Contains(strings.ToLower(text), "potd")
}

func (w *WPotdAdapter) Answer(chatId int64, txt string) []tg.Method {
	receivers := []string{strconv.FormatInt(chatId, 10)}
	methods, err := w.PushData(receivers)
	if err != nil {
		sm := tg.SendMessage[int64]{ChatId: chatId}
		sm.Text = "не получилось ("
		return []tg.Method{&sm}
	}
	return methods
}

func (w *WPotdAdapter) PushData(receivers []string) ([]tg.Method, error) {
	t, err := w.service.Potd()
	if err != nil {
		return nil, err
	}

	res := make([]tg.Method, 0)
	for _, recv := range receivers {
		sm := new(tg.SendPhoto[string])
		sm.ChatId = recv
		sm.Photo = t.Src
		sm.Text = t.Text
		sm.SendOptions = tg.SendOptions{DisableNotification: true}
		res = append(res, sm)
	}
	return res, nil
}
