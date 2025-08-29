package wiki

import (
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/bochkov/gobot/internal/lib/htmlq"
	"github.com/bochkov/gobot/internal/services"
	"github.com/microcosm-cc/bluemonday"
)

type ItdService struct {
	services.Service
	p *bluemonday.Policy
}

func NewItd() *ItdService {
	p := bluemonday.NewPolicy()
	p.AllowElements("b", "strong", "i", "em")
	p.AllowAttrs("href").OnElements("a")
	return &ItdService{p: p}
}

func (s *ItdService) Description() string {
	return "В этот день"
}

func (s *ItdService) Itd() (*ThisDay, error) {
	doc, err := htmlq.LoadURL(WikiUrl)
	if err != nil {
		return nil, err
	}

	todayNode, err := htmlquery.Query(doc, "//div[@id='main-itd']")
	if err != nil {
		return nil, err
	}

	var res ThisDay

	date, err := htmlquery.Query(todayNode, "//h2/div[2]/a")
	if err != nil {
		return nil, err
	}
	res.Date = sanitizeHtml(s.p, htmlquery.OutputHTML(date, true))

	worldDay, _ := htmlquery.Query(todayNode, "//p")
	if worldDay != nil {
		res.WorldDay = sanitizeHtml(s.p, htmlquery.OutputHTML(worldDay, true))
	}

	img, _ := htmlquery.Query(todayNode, "//figure/a/img")
	if img != nil {
		src := "https:" + htmlquery.SelectAttr(img, "src")
		res.ImgSrc = regexp.MustCompile(`\d+px`).ReplaceAllString(src, "1000px")
	}

	list, err := htmlquery.QueryAll(todayNode, "//ul/li")
	if err != nil {
		return nil, err
	}

	var sb strings.Builder
	for _, n := range list {
		sb.WriteString("\n")
		sb.WriteString(sanitizeHtml(s.p, htmlquery.OutputHTML(n, false)))
	}
	res.Text = sb.String()

	return &res, nil
}
