package services

import (
	"fmt"
	"speedstar/internal/db"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func RunScheduler() {
	// create a scheduler
	s, err := gocron.NewScheduler()
	// defer func() { _ = s.Shutdown() }()
	// defer s.Shutdown()
	if err != nil {
		// handle error
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
		// handle error
	}
	// each job has a unique id
	fmt.Println(j.ID())

	_, _ = s.NewJob(
		gocron.CronJob(
			// standard cron tab parsing
			"1 * * * *",
			false,
		),
		gocron.NewTask(
			func() {},
		),
	)
	_, _ = s.NewJob(
		gocron.CronJob(
			// optionally include seconds as the first field
			"* 1 * * * *",
			true,
		),
		gocron.NewTask(
			func() {},
		),
	)

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
