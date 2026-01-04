package scheduler

import (
	"log"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron *cron.Cron
}

var schedulerInstance *Scheduler

func InitScheduler() {
	schedulerInstance = &Scheduler{
		cron: cron.New(),
	}
	schedulerInstance.cron.Start()
	log.Println("Scheduler initialized")
}

func GetScheduler() *Scheduler {
	return schedulerInstance
}

func (s *Scheduler) AddJob(spec string, cmd func()) (cron.EntryID, error) {
	return s.cron.AddFunc(spec, cmd)
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}
