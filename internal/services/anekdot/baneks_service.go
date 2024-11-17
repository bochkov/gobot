package anekdot

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/carlmjohnson/requests"
)

const url = "https://baneks.ru"

type baneks struct {
}

func NewService() Service {
	return &baneks{}
}

func (b *baneks) Description() string {
	return "Случайный анекдот от https://baneks.ru/"
}

func (b *baneks) GetRandom() (*Anekdot, error) {
	return b.getByUrl(fmt.Sprintf("%s/random", url))
}

func (b *baneks) GetById(id int) (*Anekdot, error) {
	return b.getByUrl(fmt.Sprintf("%s/%d", url, id))
}

func (b *baneks) getByUrl(url string) (*Anekdot, error) {
	var page string
	if err := requests.
		URL(url).
		ToString(&page).
		Fetch(context.Background()); err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(page))
	if err != nil {
		return nil, err
	}

	var anek = &Anekdot{}
	article := doc.Find("article")
	id := article.Find("h2").Text()
	idx := strings.Index(id, "#")
	if idx > -1 {
		anek.Id, _ = strconv.Atoi(id[idx+1:])
	} else {
		anek.Id = -1
	}
	anek.Text = doc.Find("article p").Text()
	return anek, nil
}
