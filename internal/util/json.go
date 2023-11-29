package util

import (
	"encoding/json"
	"log"
	"net/http"
)

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
