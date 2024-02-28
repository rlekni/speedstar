package services

import (
	"fmt"
	"log"
	"os"
	"speedstar/internal/db"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type ISpeedtestScheduler interface {
	RunScheduler()
}

type SpeedtestScheduler struct {
	scheduler gocron.Scheduler
}

func NewScheduler() ISpeedtestScheduler {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Println(err)
	}
	return &SpeedtestScheduler{
		scheduler: scheduler,
	}
}

func (scheduler SpeedtestScheduler) RunScheduler() {
	cronTab := os.Getenv("SCHEDULE_CRON")
	log.Printf("Cron: %s\n", cronTab)
	cronJob, err := scheduler.scheduler.NewJob(
		gocron.CronJob(
			// standard cron tab parsing
			cronTab,
			false,
		),
		gocron.NewTask(
			func() {
				// Create new service/clients every run
				var service = NewSpeedtestService()
				service.RunSpeedtest()
			},
		),
	)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(cronJob.ID())
	log.Printf("Cron Job created: %s\n", cronJob.ID())
	scheduler.scheduler.Start()
}

func RunScheduler() {
	// create a scheduler
	s, err := gocron.NewScheduler()
	// defer func() { _ = s.Shutdown() }()
	// defer s.Shutdown()
	if err != nil {
		log.Println(err)
	}

	// add a job to the scheduler
	j, err := s.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			func(a string, b int) {
				// do things
				db.InfluxdbWrites()
				db.InfluxdbWritesAsync()
				fmt.Println("JOB RUN")
			},
			"hello",
			1,
		),
	)
	if err != nil {
		log.Println(err)
	}
	// each job has a unique id
	fmt.Println(j.ID())

	_, err = s.NewJob(
		gocron.CronJob(
			// standard cron tab parsing
			os.Getenv("SCHEDULE_CRON"),
			false,
		),
		gocron.NewTask(
			func() {}, // speed test
		),
	)
	// _, _ = s.NewJob(
	// 	gocron.CronJob(
	// 		// optionally include seconds as the first field
	// 		"* 1 * * * *",
	// 		true,
	// 	),
	// 	gocron.NewTask(
	// 		func() {},
	// 	),
	// )
	if err != nil {
		// handle error
		log.Println(err)
	}

	// start the scheduler
	s.Start()

	// block until you are ready to shut down
	// select {
	// case <-time.After(time.Minute):
	// }

	// // when you're done, shut it down
	// err = s.Shutdown()
	// if err != nil {
	// 	// handle error
	// }
}
