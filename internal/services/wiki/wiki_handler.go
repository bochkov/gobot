package wiki

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

func (h *Handler) WikiHandler(w http.ResponseWriter, req *http.Request) {
	an, err := h.Service.Today()
	util.JsonResponse(w, an, err)
}
