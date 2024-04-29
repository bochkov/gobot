package tg

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/bochkov/gobot/internal/cbr"
)

type CbrWorker struct {
	cbr.Service
}

func NewCbrWorker(s cbr.Service) *CbrWorker {
	return &CbrWorker{Service: s}
}

func (c *CbrWorker) Description() string {
	return c.Service.Description()
}

func (c *CbrWorker) IsMatch(text string) bool {
	return strings.Contains(strings.ToLower(text), "курс")
}

func (c *CbrWorker) Answer(chatId int64, txt string) []Method {
	text := strings.ToUpper(txt)
	currencies := regexp.MustCompile(`\s+`).Split(text, -1)[1:]
	if len(currencies) == 0 {
		currencies = append(currencies, "USD")
		currencies = append(currencies, "EUR")
	}
	var methods = make([]Method, 0)
	ranges := c.Service.LatestRange(context.Background(), currencies)
	if len(ranges) != 0 {
		sm := SendMessage[int64]{
			ChatId: chatId,
			Text:   fmt.Sprintf("Курс на %s", ranges[0].Value.Date.Time.Format("02.01.2006")),
		}
		methods = append(methods, &sm)
	}
	for _, r := range ranges {
		sm := SendMessage[int64]{
			ChatId: chatId,
			Text:   r.String(),
		}
		methods = append(methods, &sm)
	}
	return methods
}
