package util

import (
	"net/http"
)

type Err struct {
	Description string `json:"description"`
}

func JsonResponse(w http.ResponseWriter, model any, err error) {
	if err == nil {
		AsJson(w, http.StatusOK, model)
	} else {
		e := &Err{Description: err.Error()}
		AsJson(w, http.StatusInternalServerError, e)
	}
}
