package rutor

import (
	"github.com/bochkov/gobot/internal/services"
	"golang.org/x/net/html"
)

type Torrent struct {
	MagnetUrl string
	DirectUrl string
	Bytes     []byte
	Name      string
}

type Service interface {
	services.Service
	FetchTorrent(url string) (*Torrent, error)
}

func getAttrValue(sel *html.Node, key string) string {
	for _, a := range sel.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}
