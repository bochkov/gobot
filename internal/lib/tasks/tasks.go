package tasks

import (
	"log"

	"github.com/bochkov/gobot/internal/lib/db"
	"github.com/bochkov/gobot/internal/push"
	"github.com/go-co-op/gocron"
)

type SchedParam struct {
	Desc     string
	CronProp string
	CronDef  string
}

func Schedule(scheduler *gocron.Scheduler, pushSrv push.Service, push push.Push, param SchedParam) {
	cron := db.GetProp(param.CronProp, param.CronDef)
	_, err := scheduler.Cron(cron).Do(func() {
		log.Printf("execute %s", param.Desc)
		text := push.PushText()
		if text == "" {
			return
		}
		pushSrv.Push(text)
	})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("scheduled %s at '%s'", param.Desc, cron)
	}
}
