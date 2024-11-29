package tasks

import (
	"log/slog"

	"github.com/bochkov/gobot/internal/lib/db"
	"github.com/bochkov/gobot/internal/push"
	"github.com/go-co-op/gocron"
)

type Scheduled interface {
	Schedule(schedule *gocron.Scheduler)
}

type SchedParam struct {
	Desc       string
	CronProp   string
	CronDef    string
	Recipients []string
}

func Schedule(scheduler *gocron.Scheduler, service push.Service, param SchedParam) {
	cron := db.GetProp(param.CronProp, param.CronDef)
	_, err := scheduler.Cron(cron).Do(func() {
		slog.Info("execute tasks", "task", param.Desc)
		service.Push(param.Recipients)
	})
	if err != nil {
		slog.Warn(err.Error())
	} else {
		slog.Info("scheduled", "task", param.Desc, "at", cron)
	}
}
