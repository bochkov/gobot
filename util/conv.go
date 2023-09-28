package util

import (
	"encoding/json"
	"encoding/xml"
	"golang.org/x/net/html/charset"
	"log"
	"net/http"
	"strings"
)

type Err struct {
	Description string `json:"description"`
}

func NewErr(desc string) *Err {
	return &Err{Description: desc}
}

func FromXml(data string, model any) error {
	decoder := xml.NewDecoder(strings.NewReader(data))
	decoder.CharsetReader = charset.NewReaderLabel
	if err := decoder.Decode(model); err != nil {
		return err
	}
	return nil
}

func ToJson(model any) string {
	js, _ := json.Marshal(model)
	return string(js)
}

func FromJsonString(data string) (map[string]any, error) {
	x := map[string]any{}
	if err := json.Unmarshal([]byte(data), &x); err != nil {
		return nil, err
	}
	return x, nil
}

func FromJson(req *http.Request, model any) error {
	if err := json.NewDecoder(req.Body).Decode(&model); err != nil {
		return err
	}
	return nil
}

func AsJson(w http.ResponseWriter, status int, model any) {
	js, err := json.Marshal(model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err = w.Write(js); err != nil {
		log.Println(err)
		return
	}
}

func Resp(w http.ResponseWriter, model any, err error) {
	if err == nil {
		AsJson(w, http.StatusOK, model)
	} else {
		AsJson(w, http.StatusInternalServerError, NewErr(err.Error()))
	}
}
