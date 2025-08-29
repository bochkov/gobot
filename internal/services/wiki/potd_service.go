package wiki

import (
	"github.com/antchfx/htmlquery"
	"github.com/bochkov/gobot/internal/lib/htmlq"
	"github.com/bochkov/gobot/internal/services"
	"github.com/microcosm-cc/bluemonday"
)

type PotdService struct {
	services.Service
	p *bluemonday.Policy
}

func NewPotd() *PotdService {
	p := bluemonday.NewPolicy()
	p.AllowElements("b", "strong", "i", "em")
	p.AllowAttrs("href").OnElements("a")
	return &PotdService{p: p}
}

func (s *PotdService) Description() string {
	return "Изображение дня"
}

func (s *PotdService) Potd() (*PicOfDay, error) {
	doc, err := htmlq.LoadURL(WikiUrl)
	if err != nil {
		return nil, err
	}

	potdNode, err := htmlquery.Query(doc, "//div[@id='main-potd']")
	if err != nil {
		return nil, err
	}

	var res PicOfDay

	img, err := htmlquery.Query(potdNode, "//div[1]/div[1]/figure/a/img")
	if err != nil {
		return nil, err

	}
	res.Text = htmlquery.SelectAttr(img, "alt")
	res.Src = "https:" + htmlquery.SelectAttr(img, "src")
	return &res, nil
}
