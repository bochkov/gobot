package quote

import (
	"fmt"
	"log"
)

func (s *forismatic) PushText() string {
	quote, err := s.RandomQuote()
	if err != nil {
		log.Print(err)
		return ""
	}
	return fmt.Sprintf("Мудрость дня:\n%s", quote)
}
