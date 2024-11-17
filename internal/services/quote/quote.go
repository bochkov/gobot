package quote

import (
	"fmt"

	"github.com/bochkov/gobot/internal/services"
)

type Quote struct {
	QuoteText   string `json:"quoteText"`
	QuoteAuthor string `json:"quoteAuthor"`
	SenderName  string `json:"senderName"`
	SenderLink  string `json:"senderLink"`
	QuoteLink   string `json:"quoteLink"`
}

func (q *Quote) String() string {
	return fmt.Sprintf("%s %s", q.QuoteText, q.QuoteAuthor)
}

type Service interface {
	services.Service
	RandomQuote() (*Quote, error)
}
