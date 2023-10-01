package tasks

import (
	"log"

	"github.com/bochkov/gobot/internal/db"
	"github.com/bochkov/gobot/internal/resnyx/anekdot"
	"github.com/bochkov/gobot/internal/tg"
	"github.com/go-co-op/gocron"
)

func configureAnekdot(scheduler *gocron.Scheduler, pushServ *tg.PushService) {
	cron := db.GetProp(db.AnekdotScheduler, "0 4 * * *")
	_, err := scheduler.Cron(cron).Do(func() {
		log.Print("push anekdot")
		text := anekdot.NewBaneks().PushText()
		if text == "" {
			return
		}
		pushServ.Push(text, func(sm *tg.SendMessage[string]) {
			sm.SendOptions.DisableWebPagePreview = true
			sm.SendOptions.DisableNotification = true
		})
	})
	if err != nil {
		log.Print(err)
	} else {
		log.Printf("schedule anekdot at '%s'", cron)
	}
}

func ConfigureTasks(scheduler *gocron.Scheduler) {
	pushServ := &tg.PushService{}
	configureAnekdot(scheduler, pushServ)
}
