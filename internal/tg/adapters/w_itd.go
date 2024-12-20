package adapters

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bochkov/gobot/internal/services/wiki"
	"github.com/bochkov/gobot/internal/tg"
)

type WItdAdapter struct {
	service *wiki.ItdService
}

func NewWItdAdapter(s *wiki.ItdService) *WItdAdapter {
	return &WItdAdapter{service: s}
}

func (w *WItdAdapter) Description() string {
	return w.service.Description()
}

func (w *WItdAdapter) IsMatch(text string) bool {
	return strings.Contains(strings.ToLower(text), "today")
}

func (w *WItdAdapter) Answer(chatId int64, txt string) []tg.Method {
	receivers := []string{strconv.FormatInt(chatId, 10)}
	methods, err := w.PushData(receivers)
	if err != nil {
		sm := tg.SendMessage[int64]{ChatId: chatId}
		sm.Text = "не получилось ("
		return []tg.Method{&sm}
	}
	return methods
}

func (w *WItdAdapter) PushData(receivers []string) ([]tg.Method, error) {
	t, err := w.service.Itd()
	if err != nil {
		return nil, err
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<b>%s</b>\n", t.Date))
	if t.WorldDay != "" {
		sb.WriteString(t.WorldDay)
	}
	sb.WriteString(t.Text)

	res := make([]tg.Method, 0)
	for _, recv := range receivers {
		sm := new(tg.SendMessage[string])
		sm.ChatId = recv
		sm.Text = sb.String()
		sm.SendOptions = tg.SendOptions{ParseMode: tg.HTML, DisableNotification: true}
		sm.SendOptions.LinkPreviewOpts = &tg.LinkPreviewOptions{
			Url: t.ImgSrc, PreferLarge: true, ShowAboveText: false,
		}
		res = append(res, sm)
	}
	return res, nil
}
