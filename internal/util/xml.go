package util

import (
	"encoding/xml"
	"strings"

	"golang.org/x/net/html/charset"
)

func FromXml(data string, model any) error {
	decoder := xml.NewDecoder(strings.NewReader(data))
	decoder.CharsetReader = charset.NewReaderLabel
	if err := decoder.Decode(model); err != nil {
		return err
	}
	return nil
}
