package cbr

import (
	"net/http"
	"strings"
	"time"

	"github.com/bochkov/gobot/internal/util"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) LatestRate(w http.ResponseWriter, req *http.Request) {
	curRate, err := h.Service.LatestRate(req.Context())
	util.JsonResponse(w, curRate, err)
}

func (h *Handler) LatestRates(w http.ResponseWriter, req *http.Request) {
	currency := req.URL.Query().Get("currency")
	if currency == "" {
		currency = "USD,EUR"
	}
	ranges := h.Service.LatestRange(req.Context(), strings.Split(strings.ToUpper(currency), ","))
	util.JsonResponse(w, ranges, nil)
}

func (h *Handler) PeriodRates(w http.ResponseWriter, req *http.Request) {
	period := chi.URLParam(req, "period")
	currency := strings.ToUpper(chi.URLParam(req, "currency"))
	t := time.Now()
	switch period {
	case "month":
		r, err := h.Service.RangeOf(req.Context(), currency, t.AddDate(0, -1, 0), t.AddDate(0, 0, 1))
		util.JsonResponse(w, r, err)
	case "year":
		r, err := h.Service.RangeOf(req.Context(), currency, t.AddDate(-1, 0, 0), t.AddDate(0, 0, 1))
		util.JsonResponse(w, r, err)
	}
}
