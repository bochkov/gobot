package tg

import (
	"strings"

	"github.com/bochkov/gobot/internal/wiki"
)

type WikiWorker struct {
	wiki.Service
}

func NewWikiWorker(s wiki.Service) *WikiWorker {
	return &WikiWorker{Service: s}
}

func (w *WikiWorker) Description() string {
	return w.Service.Description()
}

func (w *WikiWorker) IsMatch(text string) bool {
	return strings.Contains(strings.ToLower(text), "today")
}

func (w *WikiWorker) Answer(chatId int64, txt string) []Method {
	today, err := wiki.NewService().Today()
	sm := SendMessage[int64]{ChatId: chatId}
	if err != nil {
		sm.Text = "не получилось ("
	}
	sm.Text = today.AsHtml()
	sm.SendOptions = SendOptions{ParseMode: HTML}
	return []Method{&sm}
}
