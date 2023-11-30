package anekdot

import (
	"fmt"
	"log/slog"
)

func (b *baneks) PushText() string {
	anek, err := b.GetRandom()
	if err != nil {
		slog.Warn("push anekdot", "err", err)
		return ""
	}
	return fmt.Sprintf("Анекдот дня:\n%s", anek.Text)
}
