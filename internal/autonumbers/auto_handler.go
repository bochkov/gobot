package autonumbers

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

func (h *Handler) AutonumbersHandler(w http.ResponseWriter, req *http.Request) {
	code := req.URL.Query().Get("code")
	if code != "" {
		region, err := h.Service.FindRegionByCode(req.Context(), code)
		util.JsonResponse(w, region, err)
		return

	}

	region := req.URL.Query().Get("region")
	if region != "" {
		codes, err := h.Service.FindRegionByName(req.Context(), region)
		util.JsonResponse(w, codes, err)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}
