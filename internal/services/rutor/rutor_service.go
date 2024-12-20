package rutor

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/carlmjohnson/requests"
)

type service struct {
}

func NewService() Service {
	return &service{}
}

func (s *service) Description() string {
	return "fetch torrents from rutor"
}

func (s *service) FetchTorrent(url string) (*Torrent, error) {
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

func (s *service) getUrls(url string) (*Torrent, error) {
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
			if strings.HasPrefix(val, "magnet:") {
				tor.MagnetUrl = val
			}
			if strings.HasPrefix(val, "//d.rutor") {
				tor.DirectUrl = "http:" + val
			}
		}
	}
	if tor.MagnetUrl == "" || tor.DirectUrl == "" {
		return nil, errors.New("no urls detected")
	}
	return tor, nil
}

func (s *service) fetch(tor *Torrent) error {
	ctx := context.Background()
	headers := http.Header{}
	if err := requests.URL(tor.DirectUrl).
		ToHeaders(headers).
		Method(http.MethodGet).
		Fetch(ctx); err != nil {
		return err
	}
	disposition := headers.Get("Content-Disposition")
	re := regexp.MustCompile(`.*filename="(?P<name>.*)"`)
	var name = re.FindStringSubmatch(disposition)[re.SubexpIndex("name")]
	if name == "" {
		name = "noname"
	}
	tor.Name = name

	buf := new(bytes.Buffer)
	if err := requests.URL(tor.DirectUrl).
		ToBytesBuffer(buf).
		Fetch(ctx); err != nil {
		return err
	}
	tor.Bytes = buf.Bytes()
	return nil
}
