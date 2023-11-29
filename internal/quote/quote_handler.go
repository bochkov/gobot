package quote

import (
	"net/http"

	"github.com/bochkov/gobot/internal/util"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) RandomQuote(w http.ResponseWriter, _ *http.Request) {
	quote, err := h.Service.RandomQuote()
	util.JsonResponse(w, quote, err)
}
