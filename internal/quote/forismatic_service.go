package quote

import (
	"context"

	"github.com/carlmjohnson/requests"
)

type forismatic struct {
}

func NewService() Service {
	return &forismatic{}
}

func (f *forismatic) RandomQuote() (*Quote, error) {
	var quote Quote
	if err := requests.
		URL("https://api.forismatic.com/api/1.0/").
		Param("method", "getQuote").
		Param("format", "json").
		Param("lang", "ru").
		ToJSON(&quote).
		Fetch(context.Background()); err != nil {
		return nil, err
	}
	return &quote, nil
}

func (f *forismatic) Description() string {
	return "Случайная цитата от forismatic.com"
}
