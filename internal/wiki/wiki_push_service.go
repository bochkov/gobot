package wiki

import (
	"log/slog"
)

func (b *today) PushText() string {
	txt, err := b.Today()
	if err != nil {
		slog.Warn("push wiki today", "err", err)
		return ""
	}
	return txt.AsHtml()
}
