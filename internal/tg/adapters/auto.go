package adapters

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/bochkov/gobot/internal/services/autonumbers"
	"github.com/bochkov/gobot/internal/tg"
)

type AutoAdapter struct {
	service autonumbers.Service
}

func NewAutoAdapter(s autonumbers.Service) AutoAdapter {
	return AutoAdapter{service: s}
}

func (a AutoAdapter) Description() string {
	return a.service.Description()
}

func (a AutoAdapter) IsMatch(text string) bool {
	match := strings.ToLower(text)
	return strings.Contains(match, "avto") ||
		strings.Contains(match, "auto") ||
		strings.Contains(match, "авто")
}

func (a AutoAdapter) Answer(chatId int64, txt string) []tg.Method {
	ctx := context.Background()
	re := regexp.MustCompile(`(?P<code>\d+)`)
	matches := re.FindAllStringSubmatch(txt, -1)
	if len(matches) == 0 {
		name := txt[strings.Index(txt, " ")+1:]
		regions, err := a.service.FindRegionByName(ctx, name)
		if err != nil {
			slog.Warn(err.Error())
		} else {
			return a.createMessages(chatId, regions)
		}
	} else {
		regions := make([]autonumbers.Region, 0)
		for _, digits := range matches {
			region, err := a.service.FindRegionByCode(ctx, digits[1])
			if err != nil {
				slog.Warn(err.Error())
			} else {
				regions = append(regions, *region)
			}
		}
		return a.createMessages(chatId, regions)
	}
	return a.createMessages(chatId, []autonumbers.Region{})
}

func (a AutoAdapter) createMessages(chatId int64, regions []autonumbers.Region) []tg.Method {
	messages := make([]tg.Method, 0)
	if len(regions) == 0 {
		sm := tg.SendMessage[int64]{
			ChatId: chatId,
			Text:   "ничего не нашел (",
		}
		messages = append(messages, &sm)
	} else {
		for _, r := range regions {
			sm := tg.SendMessage[int64]{
				ChatId: chatId,
				Text:   fmt.Sprintf("%s = %s", r.Name, strings.Join(r.Codes, ", ")),
			}
			messages = append(messages, &sm)
		}
	}
	return messages
}
