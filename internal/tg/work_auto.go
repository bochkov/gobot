package tg

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"log/slog"

	"github.com/bochkov/gobot/internal/autonumbers"
)

type AutoWorker struct {
	autonumbers.Service
}

func NewAutoWorker(s autonumbers.Service) *AutoWorker {
	return &AutoWorker{Service: s}
}

func (a *AutoWorker) Description() string {
	return a.Service.Description()
}

func (a *AutoWorker) IsMatch(text string) bool {
	match := strings.ToLower(text)
	return strings.Contains(match, "avto") ||
		strings.Contains(match, "auto") ||
		strings.Contains(match, "авто")
}

func (a *AutoWorker) Answer(chatId int64, txt string) []Method {
	ctx := context.Background()
	re := regexp.MustCompile(`(?P<code>\d+)`)
	matches := re.FindAllStringSubmatch(txt, -1)
	if len(matches) == 0 {
		name := txt[strings.Index(txt, " ")+1:]
		regions, err := a.Service.FindRegionByName(ctx, name)
		if err != nil {
			slog.Warn(err.Error())
		} else {
			return a.createMessages(chatId, regions)
		}
	} else {
		regions := make([]autonumbers.Region, 0)
		for _, digits := range matches {
			region, err := a.Service.FindRegionByCode(ctx, digits[1])
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

func (a *AutoWorker) createMessages(chatId int64, regions []autonumbers.Region) []Method {
	messages := make([]Method, 0)
	if len(regions) == 0 {
		sm := SendMessage[int64]{
			ChatId: chatId,
			Text:   "ничего не нашел (",
		}
		messages = append(messages, &sm)
	} else {
		for _, r := range regions {
			sm := SendMessage[int64]{
				ChatId: chatId,
				Text:   fmt.Sprintf("%s = %s", r.Name, strings.Join(r.Codes, ", ")),
			}
			messages = append(messages, &sm)
		}
	}
	return messages
}
