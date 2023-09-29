package main

import (
	"context"
	"fmt"
	"github.com/bochkov/gobot/internal/db"
	"github.com/bochkov/gobot/internal/resnyx/anekdot"
	"github.com/bochkov/gobot/internal/resnyx/auto"
	"github.com/bochkov/gobot/internal/resnyx/cbr"
	"github.com/bochkov/gobot/internal/resnyx/forismatic"
	"github.com/bochkov/gobot/internal/tasks"
	"github.com/bochkov/gobot/internal/tg"
	"github.com/bochkov/gobot/internal/util"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
)

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
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	db.NewPool(ctx, fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pass, host, port, name))
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
	tasks.ConfigureTasks(scheduler)
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
	scheduler.Stop()
	if err := srv.Shutdown(stopCtx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
	log.Print("app stopped")
}
