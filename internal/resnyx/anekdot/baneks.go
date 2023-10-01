package anekdot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/carlmjohnson/requests"
)

const url = "https://baneks.ru"

type Baneks struct {
}

func NewBaneks() *Baneks {
	return new(Baneks)
}

func (b *Baneks) GetRandom() (*Anek, error) {
	return b.getByUrl(fmt.Sprintf("%s/random", url))
}

func (b *Baneks) GetById(id int) (*Anek, error) {
	return b.getByUrl(fmt.Sprintf("%s/%d", url, id))
}

func (b *Baneks) getByUrl(url string) (*Anek, error) {
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

	var anek = &Anek{}
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

func (b *Baneks) PushText() string {
	anek, err := b.GetRandom()
	if err != nil {
		log.Print(err)
		return ""
	}
	return fmt.Sprintf("Анекдот дня:\n%s", anek.Text)
}
