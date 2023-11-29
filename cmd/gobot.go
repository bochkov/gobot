package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bochkov/gobot/internal/anekdot"
	"github.com/bochkov/gobot/internal/autonumbers"
	"github.com/bochkov/gobot/internal/cbr"
	"github.com/bochkov/gobot/internal/lib/db"
	"github.com/bochkov/gobot/internal/lib/router"
	"github.com/bochkov/gobot/internal/lib/tasks"
	"github.com/bochkov/gobot/internal/quote"
	"github.com/bochkov/gobot/internal/rutor"
	"github.com/bochkov/gobot/internal/tg"

	"github.com/go-co-op/gocron"
)

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
	if err := db.GetPool().Ping(ctx); err != nil {
		log.Print(err)
		os.Exit(1)
	}
	var version string
	if err := db.GetPool().QueryRow(ctx, "select version()").Scan(&version); err == nil {
		log.Print(version)
	}

	/// services
	sAnekdot := anekdot.NewService()
	sQuotes := quote.NewService()
	sTorrent := rutor.NewService()
	sAutonumbers := autonumbers.NewService(
		autonumbers.NewRepository(
			db.GetPool(),
		),
	)
	sCbr := cbr.NewService(
		cbr.NewRepository(
			db.GetPool(),
		),
	)
	sTelegram := tg.NewService(
		tg.NewAnekdotWorker(sAnekdot),
		tg.NewAutoWorker(sAutonumbers),
		tg.NewQuoteWorker(sQuotes),
		tg.NewCbrWorker(sCbr),
		tg.NewRutorWorker(sTorrent),
	)

	/// scheduler
	scheduler := gocron.NewScheduler(time.UTC)
	tasks.Schedule(scheduler, sTelegram, sAnekdot, tasks.SchedParam{
		Desc: "anekdot", CronProp: db.AnekdotScheduler, CronDef: "0 4 * * *",
	})
	scheduler.StartAsync()

	/// handlers
	handlers := &router.Handlers{
		Anekdot:  anekdot.NewHandler(sAnekdot),
		Auto:     autonumbers.NewHandler(sAutonumbers),
		Cbr:      cbr.NewHandler(sCbr),
		Quotes:   quote.NewHandler(sQuotes),
		Telegram: tg.NewHandler(sTelegram),
	}
	router := router.ConfigureRouter(handlers)
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
	db.GetPool().Close()
	log.Print("app stopped")
}
