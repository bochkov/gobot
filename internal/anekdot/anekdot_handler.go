package anekdot

import (
	"net/http"
	"strconv"

	"github.com/bochkov/gobot/internal/util"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) AnekdotHandler(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil {
		an, err := h.Service.GetRandom()
		util.JsonResponse(w, an, err)
	} else {
		an, err := h.Service.GetById(id)
		util.JsonResponse(w, an, err)
	}
}
