package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
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
	"github.com/bochkov/gobot/internal/push"
	"github.com/bochkov/gobot/internal/quote"
	"github.com/bochkov/gobot/internal/rutor"
	"github.com/bochkov/gobot/internal/tg"
	"github.com/bochkov/gobot/internal/util"
	"github.com/bochkov/gobot/internal/wiki"

	"github.com/go-co-op/gocron"
	"github.com/lmittmann/tint"
)

func main() {
	/// logging
	slog.SetDefault(
		slog.New(
			tint.NewHandler(
				os.Stdout,
				&tint.Options{
					TimeFormat: time.DateTime,
					Level:      slog.LevelDebug,
					AddSource:  true,
				},
			)))

	/// params
	flags, err := util.ParseParameters()
	if err != nil {
		flag.Usage()
		panic(err)
	}

	ctx := context.Background()
	dbcp := db.NewPool(ctx, flags.DbConnectString())
	if err := dbcp.Ping(ctx); err != nil {
		slog.Error("cannot connect to db", "err", err)
		os.Exit(1)
	}
	var version string
	if err := dbcp.QueryRow(ctx, "select version()").Scan(&version); err == nil {
		slog.Info(version)
	}

	/// services
	sWikiToday := wiki.NewService()
	sAnekdot := anekdot.NewService()
	sQuotes := quote.NewService()
	sTorrent := rutor.NewService()
	sAutonumbers := autonumbers.NewService(
		autonumbers.NewRepository(dbcp),
	)
	sCbr := cbr.NewService(
		cbr.NewRepository(dbcp),
	)
	sCbrTasks := cbr.NewTaskService(
		cbr.NewTaskRepo(dbcp),
	)
	sTelegram := tg.NewService(
		tg.NewAnekdotWorker(sAnekdot),
		tg.NewAutoWorker(sAutonumbers),
		tg.NewQuoteWorker(sQuotes),
		tg.NewCbrWorker(sCbr),
		tg.NewRutorWorker(sTorrent),
		tg.NewWikiWorker(sWikiToday),
	)

	/// scheduler
	scheduler := gocron.NewScheduler(time.UTC)
	tasks.Schedule(scheduler, sTelegram, sWikiToday, tasks.SchedParam{
		Desc: "wiki today", CronProp: db.WikiScheduler, CronDef: "0 6 * * *",
	})
	sCbrTasks.Schedule(scheduler)
	scheduler.StartAsync()

	/// handlers
	handlers := &router.Handlers{
		Anekdot:  anekdot.NewHandler(sAnekdot),
		Auto:     autonumbers.NewHandler(sAutonumbers),
		Cbr:      cbr.NewHandler(sCbr),
		Quotes:   quote.NewHandler(sQuotes),
		Wiki:     wiki.NewHandler(sWikiToday),
		Telegram: tg.NewHandler(sTelegram),
		Dev:      push.NewHandler(sTelegram, sWikiToday),
	}
	routes := router.ConfigureRouter(handlers)
	srv := &http.Server{Addr: flags.ServeAddr(), Handler: routes}

	// start
	notifyCtx, nStop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer nStop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("cannot start listener", "err", err)
		}
	}()
	slog.Info(fmt.Sprintf("app started at addr='%s'", srv.Addr))
	<-notifyCtx.Done()

	slog.Info("stopping app")
	stopCtx, cStop := context.WithTimeout(ctx, 5*time.Second)
	defer cStop()

	scheduler.Stop()
	if err := srv.Shutdown(stopCtx); err != nil {
		slog.Warn("shutdown", "err", err)
	}
	dbcp.Close()
	slog.Info("app stopped")
}
