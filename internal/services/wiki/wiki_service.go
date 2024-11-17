package wiki

import (
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/microcosm-cc/bluemonday"
)

const WikiUrl string = "https://ru.wikipedia.org/wiki/%D0%97%D0%B0%D0%B3%D0%BB%D0%B0%D0%B2%D0%BD%D0%B0%D1%8F_%D1%81%D1%82%D1%80%D0%B0%D0%BD%D0%B8%D1%86%D0%B0"

type today struct {
	sanitizer *bluemonday.Policy
}

func NewService() Service {
	p := bluemonday.NewPolicy()
	p.AllowElements("b", "strong", "i", "em")
	p.AllowAttrs("href").OnElements("a")
	return &today{sanitizer: p}
}

func (t *today) Description() string {
	return "В этот день"
}

func (t *today) sanitizeHtml(source string) string {
	html := strings.ReplaceAll(
		source, "/wiki", "https://ru.wikipedia.org/wiki",
	)
	return t.sanitizer.Sanitize(html)
}

func (t *today) Today() (*ThisDay, error) {
	doc, err := htmlquery.LoadURL(WikiUrl)
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
	res.Date = t.sanitizeHtml(htmlquery.OutputHTML(date, true))

	worldDay, _ := htmlquery.Query(todayNode, "//p")
	if worldDay != nil {
		res.WorldDay = t.sanitizeHtml(htmlquery.OutputHTML(worldDay, true))
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
		sb.WriteString(t.sanitizeHtml(htmlquery.OutputHTML(n, false)))
	}
	res.Text = sb.String()

	return &res, nil
}
