package anekdot

import (
	"fmt"
	"log"
)

func (b *baneks) PushText() string {
	anek, err := b.GetRandom()
	if err != nil {
		log.Print(err)
		return ""
	}
	return fmt.Sprintf("Анекдот дня:\n%s", anek.Text)
}
