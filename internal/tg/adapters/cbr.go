package adapters

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/bochkov/gobot/internal/services/cbr"
	"github.com/bochkov/gobot/internal/tg"
)

type CbrAdapter struct {
	service cbr.Service
}

func NewCbrAdapter(s cbr.Service) CbrAdapter {
	return CbrAdapter{service: s}
}

func (c CbrAdapter) Description() string {
	return c.service.Description()
}

func (c CbrAdapter) IsMatch(text string) bool {
	return strings.Contains(strings.ToLower(text), "курс") || strings.Contains(strings.ToLower(text), "exchange")
}

func (c CbrAdapter) Answer(chatId int64, txt string) []tg.Method {
	text := strings.ToUpper(txt)
	currencies := regexp.MustCompile(`\s+`).Split(text, -1)[1:]
	if len(currencies) == 0 {
		currencies = append(currencies, "USD")
		currencies = append(currencies, "EUR")
	}
	var methods = make([]tg.Method, 0)
	ranges := c.service.LatestRange(context.Background(), currencies)
	if len(ranges) != 0 {
		sm := tg.SendMessage[int64]{
			ChatId: chatId,
			Text:   fmt.Sprintf("Курс на %s", ranges[0].Value.Date.Time.Format("02.01.2006")),
		}
		methods = append(methods, &sm)
	}
	for _, r := range ranges {
		sm := tg.SendMessage[int64]{
			ChatId: chatId,
			Text:   r.String(),
		}
		methods = append(methods, &sm)
	}
	return methods
}
