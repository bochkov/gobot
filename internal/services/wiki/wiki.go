package wiki

import (
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

const WikiUrl string = "https://ru.wikipedia.org/wiki/%D0%97%D0%B0%D0%B3%D0%BB%D0%B0%D0%B2%D0%BD%D0%B0%D1%8F_%D1%81%D1%82%D1%80%D0%B0%D0%BD%D0%B8%D1%86%D0%B0"

type ThisDay struct {
	Date     string `json:"date"`
	WorldDay string `json:"worldday,omitempty"`
	Text     string `json:"text"`
	ImgSrc   string `json:"img,omitempty"`
}

type PicOfDay struct {
	Text string `json:"text"`
	Src  string `json:"src"`
}

func sanitizeHtml(p *bluemonday.Policy, source string) string {
	html := strings.ReplaceAll(
		source, "/wiki", "https://ru.wikipedia.org/wiki",
	)
	return p.Sanitize(html)
}
