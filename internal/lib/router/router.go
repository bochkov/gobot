package router

import (
	"time"

	"github.com/bochkov/gobot/internal/services/anekdot"
	"github.com/bochkov/gobot/internal/services/autonumbers"
	"github.com/bochkov/gobot/internal/services/cbr"
	"github.com/bochkov/gobot/internal/services/dev"
	"github.com/bochkov/gobot/internal/services/quote"
	"github.com/bochkov/gobot/internal/services/wiki"
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
	Wiki     *wiki.Handler
	Telegram *tg.Handler
	Dev      *dev.Handler
}

func ConfigureRouter(h *Handlers, dev bool) *chi.Mux {
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

	r.Route("/wiki", func(r chi.Router) {
		r.Get("/today", h.Wiki.Itd)
		r.Get("/potd", h.Wiki.Potd)
	})

	if dev {
		r.Get("/bot/push", h.Dev.DevHandler)
	}

	return r
}
