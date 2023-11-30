package router

import (
	"time"

	"github.com/bochkov/gobot/internal/anekdot"
	"github.com/bochkov/gobot/internal/autonumbers"
	"github.com/bochkov/gobot/internal/cbr"
	"github.com/bochkov/gobot/internal/quote"
	"github.com/bochkov/gobot/internal/tg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var r *chi.Mux

type Handlers struct {
	Anekdot  *anekdot.Handler
	Auto     *autonumbers.Handler
	Cbr      *cbr.Handler
	Quotes   *quote.Handler
	Telegram *tg.Handler
}

func ConfigureRouter(h *Handlers) *chi.Mux {
	r = chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Post("/bot/{token}", h.Telegram.BotHandler)

	r.Get("/quote", h.Quotes.RandomQuote)

	r.Get("/anekdot", h.Anekdot.AnekdotHandler)
	r.Get("/auto", h.Auto.AutonumbersHandler)

	r.Get("/cbr/latest/all", h.Cbr.LatestRate)
	r.Get("/cbr/latest", h.Cbr.LatestRates)
	r.Get("/cbr/{period:month|year}/{currency}", h.Cbr.PeriodRates)

	return r
}
