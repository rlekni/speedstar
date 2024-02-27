package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/showwin/speedtest-go/speedtest"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	go runScheduler()

	for {
		time.Sleep(1000)
	}
}

var (
	dbUrl    = os.Getenv("INFLUXDB_URL")
	dbToken  = os.Getenv("INFLUXDB_TOKEN")
	dbOrg    = os.Getenv("INFLUXDB_ORG")
	dbBucket = os.Getenv("INFLUXDB_BUCKET")
)

func influxdbWrites() {
	// Create a new client using an InfluxDB server base URL and an authentication token
	client := influxdb2.NewClient(dbUrl, dbToken)
	// Ensures background processes finishes
	defer client.Close()
	// Use blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking(dbOrg, dbBucket)

	// Create point using full params constructor
	p := influxdb2.NewPoint("stat",
		map[string]string{"unit": "temperature"},
		map[string]interface{}{"avg": 24.5, "max": 45.0},
		time.Now())
	// write point immediately
	writeAPI.WritePoint(context.Background(), p)
	// Create point using fluent style
	p = influxdb2.NewPointWithMeasurement("stat").
		AddTag("unit", "temperature").
		AddField("avg", 23.2).
		AddField("max", 45.0).
		SetTime(time.Now())
	err := writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		panic(err)
	}
	// Or write directly line protocol
	line := fmt.Sprintf("stat,unit=temperature avg=%f,max=%f", 23.5, 45.0)
	err = writeAPI.WriteRecord(context.Background(), line)
	if err != nil {
		panic(err)
	}
}

func runScheduler() {
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

func testSpeed() {
	fmt.Println("Hi")
	var speedtestClient = speedtest.New()
	// Get user's network information
	user, _ := speedtestClient.FetchUserInfo()
	fmt.Printf("ISP %s\n", user.Isp)

	serverList, _ := speedtestClient.FetchServers()
	targets, _ := serverList.FindServer([]int{})
	for _, s := range targets {
		fmt.Printf("Server: %s", s.Name)
		// Please make sure your host can access this test server,
		// otherwise you will get an error.
		// It is recommended to replace a server at this time
		s.PingTest(nil)
		s.DownloadTest()
		s.UploadTest()
		fmt.Printf("Server: %s; Latency: %s, Download: %f, Upload: %f\n", s.Name, s.Latency, s.DLSpeed, s.ULSpeed)
		s.Context.Reset() // reset counter
	}
}
