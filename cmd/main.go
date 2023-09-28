package main

import (
	"context"
	"github.com/bochkov/gobot/db"
	"github.com/bochkov/gobot/internal/anekdot"
	"github.com/bochkov/gobot/internal/auto"
	"github.com/bochkov/gobot/internal/cbr"
	"github.com/bochkov/gobot/internal/forismatic"
	"github.com/bochkov/gobot/internal/tg"
	"github.com/bochkov/gobot/util"
	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func configureTasks(scheduler *gocron.Scheduler) {
	pushServ := &tg.PushService{}
	if _, err := scheduler.Cron("0 6 * * *").Do(func() {
		log.Print("push anekdot")
		text := anekdot.NewBaneks().PushText()
		if text == "" {
			return
		}
		pushServ.Push(text, func(sm *tg.SendMessage[string]) {
			sm.SendOptions.DisableWebPagePreview = true
			sm.SendOptions.DisableNotification = true
		})
	}); err != nil {
		log.Print(err)
	}
}

func configureController(router *mux.Router) {
	router.HandleFunc("/bot/{token}", tg.BotHandler).Methods(http.MethodPost)
	router.HandleFunc("/cite", forismatic.CiteHandler).Methods(http.MethodGet)
	router.HandleFunc("/anekdot", anekdot.AnekHandler).Methods(http.MethodGet)
	router.HandleFunc("/auto/forCode", auto.CodesHandler).Methods(http.MethodGet)
	router.HandleFunc("/auto/forRegion", auto.RegionsHandler).Methods(http.MethodGet)
	router.HandleFunc("/cbr/latest/all", cbr.LatestRate).Methods(http.MethodGet)
	router.HandleFunc("/cbr/latest", cbr.LatestRates).Methods(http.MethodGet)
	router.HandleFunc("/cbr/{period:month|year}/{currency}", cbr.PeriodRates).Methods(http.MethodGet)
}

func main() {
	ctx := context.Background()
	/// logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	/// db
	db.NewPool(ctx, "postgres://resnyx:resnyx@10.10.10.10:5432/resnyx")
	defer db.GetPool().Close()
	if err := db.GetPool().Ping(ctx); err != nil {
		log.Print(err)
		os.Exit(1)
	}
	var version string
	if err := db.GetPool().QueryRow(context.Background(), "select version()").Scan(&version); err == nil {
		log.Print(version)
	}

	/// scheduler
	scheduler := gocron.NewScheduler(time.UTC)
	configureTasks(scheduler)
	scheduler.StartAsync()

	/// http server
	router := mux.NewRouter()
	configureController(router)
	router.Use(util.LogMiddleware)
	srv := &http.Server{Addr: ":5000", Handler: router}

	// start
	notifyCtx, nStop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer nStop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	log.Print("app started")

	<-notifyCtx.Done()
	log.Print("stopping app")
	stopCtx, cStop := context.WithTimeout(ctx, 5*time.Second)
	defer cStop()
	if err := srv.Shutdown(stopCtx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
	log.Print("app stopped")
}
