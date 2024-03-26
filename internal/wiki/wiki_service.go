package wiki

import (
	"fmt"
	"strings"

	"github.com/antchfx/htmlquery"
)

type today struct {
}

func NewService() Service {
	return &today{}
}

func (t *today) Description() string {
	return "В этот день"
}

func (t *today) Today() (string, error) {
	doc, err := htmlquery.LoadURL("https://ru.wikipedia.org/wiki/%D0%97%D0%B0%D0%B3%D0%BB%D0%B0%D0%B2%D0%BD%D0%B0%D1%8F_%D1%81%D1%82%D1%80%D0%B0%D0%BD%D0%B8%D1%86%D0%B0")
	if err != nil {
		return "", err
	}

	todayNode, err := htmlquery.Query(doc, "//div[@id='main-itd']")
	if err != nil {
		return "", err
	}

	date, err := htmlquery.Query(todayNode, "//h2/span[2]/div[2]/a")
	if err != nil {
		return "", err
	}

	list, err := htmlquery.QueryAll(todayNode, "//ul/li")
	if err != nil {
		return "", err
	}

	img, err := htmlquery.Query(todayNode, "//figure/a/img")
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<h2>%s</h2>", htmlquery.InnerText(date)))

	if img != nil {
		sb.WriteString(
			strings.ReplaceAll(
				htmlquery.OutputHTML(img, true),
				"//",
				"https://",
			),
		)
	}

	sb.WriteString("<ul>")
	for _, n := range list {
		html := strings.ReplaceAll(
			htmlquery.OutputHTML(n, false),
			"/wiki",
			"https://ru.wikipedia.org/wiki",
		)
		sb.WriteString(fmt.Sprintf("<li>%s</li>", html))
	}
	sb.WriteString("</ul>")

	return sb.String(), nil
}
