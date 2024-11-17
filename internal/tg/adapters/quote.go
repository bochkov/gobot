package adapters

import (
	"fmt"
	"strings"

	"github.com/bochkov/gobot/internal/services/quote"
	"github.com/bochkov/gobot/internal/tg"
)

type QuoteAdapter struct {
	*tg.AnswerWithPushAdapter
	service quote.Service
}

func NewQuoteAdapter(s quote.Service) QuoteAdapter {
	return QuoteAdapter{service: s}
}

func (q QuoteAdapter) Description() string {
	return q.service.Description()
}

func (q QuoteAdapter) IsMatch(text string) bool {
	return strings.Contains(strings.ToLower(text), "цитат") || strings.Contains(strings.ToLower(text), "quote")
}

func (q QuoteAdapter) PushData(receivers []string) ([]tg.Method, error) {
	quote, err := q.service.RandomQuote()
	if err != nil {
		return nil, err
	}
	res := make([]tg.Method, 0)
	for _, recv := range receivers {
		sm := new(tg.SendMessage[string])
		sm.ChatId = recv
		sm.Text = fmt.Sprintf("Мудрость дня:\n%s", quote)
		res = append(res, sm)
	}
	return res, nil
}
