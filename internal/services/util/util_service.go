package util

import (
	"context"
	"log/slog"

	"github.com/bochkov/gobot/internal/lib/db"
	"github.com/bochkov/gobot/internal/lib/tasks"
	"github.com/bochkov/gobot/internal/repo"
	"github.com/bochkov/gobot/internal/tg"
	"github.com/go-co-op/gocron/v2"
)

type taskService struct {
	db repo.TmpMsgDao
	tg tg.Service
}

func NewTaskService(db repo.TmpMsgDao, tg tg.Service) tasks.Scheduled {
	return &taskService{db: db, tg: tg}
}

func (ts *taskService) Schedule(scheduler gocron.Scheduler) {
	var err error
	ctx := context.Background()

	// At 00:00.
	_, err = scheduler.NewJob(
		gocron.CronJob("* * * * *", false),
		gocron.NewTask(ts.doMaitenanceJob),
		gocron.WithContext(ctx),
	)
	if err != nil {
		slog.Error(err.Error())
	}
}

func (ts *taskService) doMaitenanceJob(c context.Context) {
	token := db.GetProp(db.TgBotTokenKey, "")
	if token == "" {
		slog.Warn("no token specified")
		return
	}

	rec, err := ts.db.GetAll(c)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	for _, r := range rec {
		dm := tg.DeleteMessage[int64]{ChatId: r.ChatId, MessageId: r.MsgId}
		_, err := ts.tg.Exec(&dm, token)
		if err != nil {
			slog.Warn(err.Error())
		} else {
			delErr := ts.db.Delete(c, r.Id)
			if delErr != nil {
				slog.Error(delErr.Error())
			}
		}
	}
}
