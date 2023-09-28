package rutor

import (
	"bytes"
	"context"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/carlmjohnson/requests"
	"golang.org/x/net/html"
	"net/http"
	"regexp"
	"strings"
)

type Torrent struct {
	MagnetUrl string
	DirectUrl string
	Bytes     []byte
	Name      string
}

type Service struct {
}

func NewService() *Service {
	return new(Service)
}

func getAttrValue(sel *html.Node, key string) string {
	for _, a := range sel.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func (s *Service) getUrls(url string) (*Torrent, error) {
	match, err := regexp.Match("https?://.*", []byte(url))
	if err != nil || !match {
		return nil, err
	}

	ctx := context.Background()
	var data string
	if err = requests.URL(url).ToString(&data).Fetch(ctx); err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		return nil, err
	}

	hrefs := doc.Find("#download a")
	tor := new(Torrent)
	for _, href := range hrefs.Nodes {
		val := getAttrValue(href, "href")
		if val != "" {
			if strings.HasPrefix("magnet://", val) {
				tor.MagnetUrl = val
			}
			if strings.HasPrefix("http", val) {
				tor.DirectUrl = val
			}
		}
	}
	if tor.MagnetUrl == "" || tor.DirectUrl == "" {
		return nil, errors.New("no urls detected")
	}
	return tor, nil
}

func (s *Service) fetch(tor *Torrent) error {
	ctx := context.Background()
	buf := new(bytes.Buffer)
	headers := http.Header{}
	if err := requests.URL(tor.DirectUrl).ToBytesBuffer(buf).ToHeaders(headers).Fetch(ctx); err != nil {
		return err
	}
	disposition := headers.Get("Content-Disposition")
	re := regexp.MustCompile(".* filename=\"(?P<name>.*)\"")
	var name = re.FindStringSubmatch(disposition)[re.SubexpIndex("name")]
	if name == "" {
		name = "noname"
	}
	tor.Name = name
	tor.Bytes = buf.Bytes()
	return nil
}

func (s *Service) GetTorrent(url string) (*Torrent, error) {
	tor, err := s.getUrls(url)
	if err != nil {
		return nil, err
	}
	err = s.fetch(tor)
	if err != nil {
		return nil, err
	}
	return tor, nil
}
