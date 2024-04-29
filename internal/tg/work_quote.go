package tg

import (
	"strings"

	"github.com/bochkov/gobot/internal/quote"
)

type QuoteWorker struct {
	quote.Service
}

func NewQuoteWorker(s quote.Service) *QuoteWorker {
	return &QuoteWorker{Service: s}
}

func (q *QuoteWorker) Description() string {
	return q.Service.Description()
}

func (q *QuoteWorker) IsMatch(text string) bool {
	return strings.Contains(strings.ToLower(text), "цитат") || strings.Contains(strings.ToLower(text), "quote")
}

func (q *QuoteWorker) Answer(chatId int64, txt string) []Method {
	cite, err := quote.NewService().RandomQuote()
	sm := SendMessage[int64]{ChatId: chatId}
	if err != nil {
		sm.Text = "не получилось ("
	}
	sm.Text = cite.String()
	return []Method{&sm}
}
