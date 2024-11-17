package dev

import (
	"net/http"

	"github.com/bochkov/gobot/internal/push"
)

type Handler struct {
	push.Service
}

func NewHandler(service push.Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) DevHandler(w http.ResponseWriter, req *http.Request) {
	h.Service.Push()
}
