package services

import (
	"log"
	"os"
	"speedstar/internal/db"

	"github.com/go-co-op/gocron/v2"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var (
	dbUrl   = os.Getenv("INFLUXDB_URL")
	dbToken = os.Getenv("INFLUXDB_TOKEN")
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
				log.Printf("Connecting to: %s\n", dbUrl)
				dbClient := influxdb2.NewClient(dbUrl, dbToken)
				// Ensures background processes finishes
				defer dbClient.Close()

				repository := db.NewSpeedtestRepository(dbClient)
				var service = NewSpeedtestService(repository)
				service.RunSpeedtest()
			},
		),
	)
	if err != nil {
		log.Println(err)
	}

	log.Printf("Cron Job created: %s\n", cronJob.ID())
	scheduler.scheduler.Start()
}
