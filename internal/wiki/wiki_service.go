package wiki

import (
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"
)

const WikiUrl string = "https://ru.wikipedia.org/wiki/%D0%97%D0%B0%D0%B3%D0%BB%D0%B0%D0%B2%D0%BD%D0%B0%D1%8F_%D1%81%D1%82%D1%80%D0%B0%D0%BD%D0%B8%D1%86%D0%B0"

type today struct {
}

func NewService() Service {
	return &today{}
}

func (t *today) Description() string {
	return "В этот день"
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

	date, err := htmlquery.Query(todayNode, "//h2/span[2]/div[2]/a")
	if err != nil {
		return nil, err
	}

	list, err := htmlquery.QueryAll(todayNode, "//ul/li")
	if err != nil {
		return nil, err
	}

	img, err := htmlquery.Query(todayNode, "//figure/a")
	if err != nil {
		return nil, err
	}

	var res ThisDay
	res.Date = strings.ReplaceAll(
		htmlquery.OutputHTML(date, true),
		"/wiki",
		"https://ru.wikipedia.org/wiki",
	)
	res.ImgSrc = strings.ReplaceAll(
		htmlquery.SelectAttr(img, "href"),
		"/wiki",
		"https://ru.wikipedia.org/wiki",
	)

	var sb strings.Builder
	for _, n := range list {
		sb.WriteString("\n")
		html := strings.ReplaceAll(
			htmlquery.OutputHTML(n, false),
			"/wiki",
			"https://ru.wikipedia.org/wiki",
		)
		var re = regexp.MustCompile(`title=".*?"`)
		html = re.ReplaceAllString(html, "")
		sb.WriteString(html)
	}
	res.Text = sb.String()

	return &res, nil
}
