package cbr

import (
	"context"
	"log/slog"

	"github.com/bochkov/gobot/internal/lib/tasks"
	"github.com/bochkov/gobot/internal/util"
	"github.com/carlmjohnson/requests"
	"github.com/go-co-op/gocron/v2"
)

type taskService struct {
	db TaskRepository
}

func NewTaskService(db TaskRepository) tasks.Scheduled {
	return &taskService{db: db}
}

func (ts *taskService) Schedule(scheduler gocron.Scheduler) {
	var err error
	ctx := context.Background()

	// At 12:00.
	_, err = scheduler.NewJob(
		gocron.CronJob("0 12 * * *", false),
		gocron.NewTask(ts.fetchCurrencies),
		gocron.WithContext(ctx),
	)
	if err != nil {
		slog.Error(err.Error())
	}

	// At 12:00 on every day-of-week from Monday through Friday
	_, err = scheduler.NewJob(
		gocron.CronJob("0 12 * * 1-5", false),
		gocron.NewTask(ts.fetchRates),
		gocron.WithContext(ctx),
	)
	if err != nil {
		slog.Error(err.Error())
	}

	go func() {
		ts.fetchCurrencies(ctx)
		ts.fetchRates(ctx)
	}()
}

func (ts *taskService) fetchRates(c context.Context) {
	slog.Debug("currency rates: start")
	var data string
	if err := requests.URL(DailyUrl).
		UserAgent("curl/8.0.1").
		ToString(&data).
		Fetch(c); err != nil {
		slog.Warn(err.Error())
		return
	}

	var currRate CurrRate
	if err := util.FromXml(data, &currRate); err != nil {
		slog.Warn(err.Error())
		return
	}
	slog.Debug("currency rates: fetched")

	ts.db.SaveCurrRate(c, currRate)
	slog.Debug("currency rates: saved")
}

func (ts *taskService) fetchCurrencies(c context.Context) {
	slog.Debug("currencies: start")
	var data string
	if err := requests.URL(CurrencyUrl).
		UserAgent("curl/8.0.1").
		ToString(&data).
		Fetch(c); err != nil {
		slog.Warn(err.Error())
		return
	}

	var currency Currency
	if err := util.FromXml(data, &currency); err != nil {
		slog.Warn(err.Error())
		return
	}
	slog.Debug("currencies: fetched")

	if err := ts.db.TruncCurrencyItems(c); err != nil {
		slog.Warn(err.Error())
		return
	}

	ts.db.SaveCurrency(c, currency)
	slog.Debug("currencies: saved")
}
