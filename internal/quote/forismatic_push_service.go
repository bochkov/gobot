package quote

import (
	"fmt"
	"log/slog"
)

func (f *forismatic) PushText() string {
	quote, err := f.RandomQuote()
	if err != nil {
		slog.Warn(err.Error())
		return ""
	}
	return fmt.Sprintf("Мудрость дня:\n%s", quote)
}
