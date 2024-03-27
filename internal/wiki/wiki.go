package wiki

import (
	"fmt"
	"strings"

	"github.com/bochkov/gobot/internal/push"
)

type ThisDay struct {
	Date     string `json:"date"`
	WorldDay string `json:"worldday,omitempty"`
	Text     string `json:"text"`
	ImgSrc   string `json:"img,omitempty"`
}

type Service interface {
	push.Push
	Today() (*ThisDay, error)
	Description() string
}

func (t *ThisDay) AsHtml() string {
	var sb strings.Builder

	if t.ImgSrc != "" {
		// https://stackoverflow.com/questions/38685619/how-to-send-an-embedded-image-along-with-text-in-a-message-via-telegram-bot-api/43705283#43705283
		sb.WriteString(fmt.Sprintf("<a href='%s'>&#8205;</a>", t.ImgSrc))
	}
	sb.WriteString(fmt.Sprintf("<b>%s</b>\n", t.Date))
	if t.WorldDay != "" {
		sb.WriteString(fmt.Sprintf("%s\n", t.WorldDay))
	}
	sb.WriteString(t.Text)

	return sb.String()
}
