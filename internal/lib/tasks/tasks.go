package tasks

import (
	"fmt"
	"log/slog"

	"github.com/bochkov/gobot/internal/lib/db"
	"github.com/bochkov/gobot/internal/push"
	"github.com/go-co-op/gocron"
)

type Scheduled interface {
	Schedule(schedule *gocron.Scheduler)
}

type SchedParam struct {
	Desc     string
	CronProp string
	CronDef  string
}

func Schedule(scheduler *gocron.Scheduler, pushSrv push.Service, push push.Push, param SchedParam) {
	cron := db.GetProp(param.CronProp, param.CronDef)
	_, err := scheduler.Cron(cron).Do(func() {
		slog.Info("execute tasks", "task", param.Desc)
		text := push.PushText()
		if text == "" {
			return
		}
		pushSrv.Push(text)
	})
	if err != nil {
		slog.Warn(err.Error())
	} else {
		slog.Info(fmt.Sprintf("scheduled %s at '%s'", param.Desc, cron))
	}
}
