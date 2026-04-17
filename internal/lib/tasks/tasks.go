package tasks

import (
	"log/slog"
	"strings"

	"github.com/bochkov/gobot/internal/lib/db"
	"github.com/bochkov/gobot/internal/push"
	"github.com/go-co-op/gocron/v2"
)

type Scheduled interface {
	Schedule(schedule gocron.Scheduler)
}

type SchedParam struct {
	Desc         string
	CronProp     string
	RecvProp     string
	TemporaryMsg bool
}

func Schedule(scheduler gocron.Scheduler, service push.Service, param SchedParam) {
	cron := db.GetProp(param.CronProp, "* * 1 * *")
	_, err := scheduler.NewJob(
		gocron.CronJob(cron, false),
		gocron.NewTask(func() {
			recv := db.GetProp(param.RecvProp, "")
			if recv == "" {
				slog.Info("task not executed because recipients is empty", "task", param.Desc)
				return
			}
			slog.Info("execute", "task", param.Desc)
			service.Push(strings.Split(recv, ";"))
		}),
	)
	if err != nil {
		slog.Warn(err.Error())
	} else {
		slog.Info("scheduled", "task", param.Desc, "at", cron)
	}
}
