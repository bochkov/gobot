package quote

import (
	"fmt"
	"log"
)

func (f *forismatic) PushText() string {
	quote, err := f.RandomQuote()
	if err != nil {
		log.Print(err)
		return ""
	}
	return fmt.Sprintf("Мудрость дня:\n%s", quote)
}
