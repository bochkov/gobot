package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/bochkov/gobot/internal/lib/db"
	"github.com/bochkov/gobot/internal/lib/router"
	"github.com/bochkov/gobot/internal/lib/tasks"
	"github.com/bochkov/gobot/internal/services/anekdot"
	"github.com/bochkov/gobot/internal/services/autonumbers"
	"github.com/bochkov/gobot/internal/services/cbr"
	"github.com/bochkov/gobot/internal/services/dev"
	"github.com/bochkov/gobot/internal/services/quote"
	"github.com/bochkov/gobot/internal/services/rutor"
	"github.com/bochkov/gobot/internal/services/wiki"
	"github.com/bochkov/gobot/internal/tg"
	"github.com/bochkov/gobot/internal/tg/adapters"

	"github.com/go-co-op/gocron"
	"github.com/lmittmann/tint"
)

type opts struct {
	dbHost   string
	dbPort   int
	dbName   string
	dbUser   string
	dbPasswd string
	port     int
	dev      bool
}

func (o *opts) dbConnectString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		o.dbUser, o.dbPasswd, o.dbHost, o.dbPort, o.dbName)
}

func (o *opts) serveAddr() string {
	return fmt.Sprintf(":%d", o.port)
}

func (o *opts) isOk() bool {
	return o.dbHost != "" && o.dbPort != 0 && o.dbName != "" && o.dbUser != "" && o.dbPasswd != ""
}

func obtainFromFlag(o *opts) {
	flag.StringVar(&o.dbHost, "dbhost", "", "database host")
	flag.IntVar(&o.dbPort, "dbport", 0, "database port")
	flag.StringVar(&o.dbName, "dbname", "", "database name")
	flag.StringVar(&o.dbUser, "dbuser", "", "database user login")
	flag.StringVar(&o.dbPasswd, "dbpassword", "", "database user password")
	flag.IntVar(&o.port, "port", 5000, "server port")
	flag.BoolVar(&o.dev, "dev", false, "enable dev endpoints")
	flag.Parse()
}

func obtainFromEnv(o *opts) {
	o.dbHost = os.Getenv("DB_HOST")
	o.dbPort, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	o.dbUser = os.Getenv("DB_USER")
	o.dbPasswd = os.Getenv("DB_PASSWORD")
	o.dbName = os.Getenv("DB_NAME")
}

func parseParameters() (*opts, error) {
	var opts opts
	obtainFromFlag(&opts)
	if opts.isOk() {
		return &opts, nil
	}
	obtainFromEnv(&opts)
	if opts.isOk() {
		return &opts, nil
	}
	return nil, errors.New("cannot parse parameters")
}

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
	opts, err := parseParameters()
	if err != nil {
		flag.Usage()
		panic(err)
	}

	ctx := context.Background()
	dbcp := db.NewPool(ctx, opts.dbConnectString())
	if err := dbcp.Ping(ctx); err != nil {
		slog.Error("cannot connect to db", "err", err)
		os.Exit(1)
	}
	var version string
	if err := dbcp.QueryRow(ctx, "select version()").Scan(&version); err == nil {
		slog.Info(version)
	}

	/// services
	sItd := wiki.NewItd()
	sPotd := wiki.NewPotd()
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
	sTelegram := tg.NewAnswerService(
		adapters.NewAnekdotAdapter(sAnekdot),
		adapters.NewAutoAdapter(sAutonumbers),
		adapters.NewCbrAdapter(sCbr),
		adapters.NewQuoteAdapter(sQuotes),
		adapters.NewRutorAdapter(sTorrent),
		adapters.NewWItdAdapter(sItd),
		adapters.NewWPotdAdapter(sPotd),
	)

	/// scheduler
	scheduler := gocron.NewScheduler(time.UTC)
	tasks.Schedule(scheduler,
		tg.NewPushService(adapters.NewWItdAdapter(sItd)),
		tasks.SchedParam{
			Desc:     "wiki today",
			CronProp: db.WikiScheduler,
			CronDef:  "* * * * *",
			RecvProp: db.ChatAutoSend,
		})
	tasks.Schedule(scheduler,
		tg.NewPushService(adapters.NewWPotdAdapter(sPotd)),
		tasks.SchedParam{
			Desc:     "wiki pic of the day",
			CronProp: db.WikiScheduler,
			CronDef:  "* * * * *",
			RecvProp: db.ChatIdKey,
		})
	sCbrTasks.Schedule(scheduler)
	scheduler.StartAsync()

	/// handlers
	handlers := &router.Handlers{
		Anekdot:  anekdot.NewHandler(sAnekdot),
		Auto:     autonumbers.NewHandler(sAutonumbers),
		Cbr:      cbr.NewHandler(sCbr),
		Quotes:   quote.NewHandler(sQuotes),
		Wiki:     wiki.NewHandler(sItd, sPotd),
		Telegram: tg.NewHandler(sTelegram),
		Dev:      dev.NewHandler(tg.NewPushService(adapters.NewWPotdAdapter(sPotd))),
	}
	routes := router.ConfigureRouter(handlers, opts.dev)
	srv := &http.Server{Addr: opts.serveAddr(), Handler: routes}

	// start
	notifyCtx, nStop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer nStop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("cannot start listener", "err", err)
		}
	}()
	slog.Info("app started", "addr", srv.Addr)
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
