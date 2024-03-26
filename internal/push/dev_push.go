package push

import (
	"net/http"
)

type DevHandler struct {
	Service
	Push
}

func NewHandler(s Service, push Push) *DevHandler {
	return &DevHandler{Service: s, Push: push}
}

func (h *DevHandler) DevHandler(w http.ResponseWriter, req *http.Request) {
	text := h.Push.PushText()
	if text == "" {
		return
	}
	h.Service.Push(text)
}
