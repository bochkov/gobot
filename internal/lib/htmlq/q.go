package htmlq

import (
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

// LoadURL loads the HTML document from the specified URL. Default enabling gzip on a HTTP request.
// this func from package htmlquery with specify user-agent
func LoadURL(url string) (*html.Node, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Set User-Agent
	req.Header.Add("User-Agent", "Golang_Resnyx_Bot/1.0")
	// Enable gzip compression.
	req.Header.Add("Accept-Encoding", "gzip")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	var reader io.ReadCloser

	defer func() {
		if reader != nil {
			reader.Close()
		}
	}()
	encoding := resp.Header.Get("Content-Encoding")
	switch encoding {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
	case "deflate":
		reader, err = zlib.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
	case "":
		reader = resp.Body
	default:
		return nil, fmt.Errorf("%s compression is not support", encoding)
	}

	r, err := charset.NewReader(reader, resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	return html.Parse(r)
}
