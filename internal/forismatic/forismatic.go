package forismatic

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/bochkov/gobot/util"
	"github.com/carlmjohnson/requests"
)

type Cite struct {
	QuoteText   string `json:"quoteText"`
	QuoteAuthor string `json:"quoteAuthor"`
	SenderName  string `json:"senderName"`
	SenderLink  string `json:"senderLink"`
	QuoteLink   string `json:"QuoteLink"`
}

func (c *Cite) String() string {
	return fmt.Sprintf("%s %s", c.QuoteText, c.QuoteAuthor)
}

type CiteService struct {
}

func NewCiteService() *CiteService {
	return new(CiteService)
}

func (f *CiteService) GetQuote() (*Cite, error) {
	var cite Cite
	if err := requests.
		URL("https://api.forismatic.com/api/1.0/").
		Param("method", "getQuote").
		Param("format", "json").
		Param("lang", "ru").
		ToJSON(&cite).
		Fetch(context.Background()); err != nil {
		return nil, err
	}
	return &cite, nil
}

func (f *CiteService) PushText() string {
	cite, err := f.GetQuote()
	if err != nil {
		log.Print(err)
		return ""
	}
	return fmt.Sprintf("Мудрость дня:\n%s", cite)
}

func CiteHandler(w http.ResponseWriter, _ *http.Request) {
	cite, err := NewCiteService().GetQuote()
	util.Resp(w, cite, err)
}
