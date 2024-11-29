package wiki

import (
	"net/http"

	"github.com/bochkov/gobot/internal/util"
)

type Handler struct {
	itd  *ItdService
	potd *PotdService
}

func NewHandler(itd *ItdService, potd *PotdService) *Handler {
	return &Handler{itd: itd, potd: potd}
}

func (h *Handler) Itd(w http.ResponseWriter, r *http.Request) {
	data, err := h.itd.Itd()
	util.JsonResponse(w, data, err)
}

func (h *Handler) Potd(w http.ResponseWriter, r *http.Request) {
	data, err := h.potd.Potd()
	util.JsonResponse(w, data, err)
}
