package wiki

import (
	"fmt"

	"github.com/bochkov/gobot/internal/push"
)

type ThisDay struct {
	ImgSrc string `json:"img,omitempty"`
	Date   string `json:"date"`
	Text   string `json:"text"`
}

type Service interface {
	push.Push
	Today() (*ThisDay, error)
	Description() string
}

func (t *ThisDay) AsHtml() string {
	// https://stackoverflow.com/questions/38685619/how-to-send-an-embedded-image-along-with-text-in-a-message-via-telegram-bot-api/43705283#43705283
	return fmt.Sprintf("<a href='%s'>&#8205;</a><b>%s</b>\n%s", t.ImgSrc, t.Date, t.Text)
}
