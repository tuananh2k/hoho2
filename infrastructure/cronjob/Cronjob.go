package cronjob

import (
	"github.com/robfig/cron/v3"
)

type Job struct {
	Config string `json:"config"`
}

type JobInterface interface {
	InitJob() *cron.Cron
}

var c *cron.Cron

func InitJob() {
	if c != nil {
		c.Stop()
		c = nil
	}
	job := new(Job)
	job.Config = "*/5 * * * *"
	job.StartJob()
}

func (job Job) StartJob() *cron.Cron {
	c = cron.New()
	c.AddFunc("config", func() {
		// handle call back cronjob
	})
	// c.Start() start cron
	return c
}
