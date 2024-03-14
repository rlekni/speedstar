package services

import (
	"log"
	"speedstar/internal/db"
	"speedstar/internal/types"
	"time"

	"github.com/showwin/speedtest-go/speedtest"
)

type ISpeedtestService interface {
	RunSpeedtest()
}

type SpeedtestService struct {
	client     *speedtest.Speedtest
	repository db.ISpeedtestRepository
}

func NewSpeedtestService(repo db.ISpeedtestRepository) ISpeedtestService {
	var speedtestClient = speedtest.New()
	return &SpeedtestService{
		client:     speedtestClient,
		repository: repo,
	}
}

func (service SpeedtestService) RunSpeedtest() {
	user, err := service.client.FetchUserInfo()
	if err != nil {
		log.Println(err)
	}
	log.Printf("ISP %s\n", user.Isp)

	serverList, err := service.client.FetchServers()
	if err != nil {
		log.Println(err)
	}
	targets, err := serverList.FindServer([]int{})
	if err != nil {
		log.Println(err)
	}
	for _, server := range targets {
		log.Printf("Server: %s", server.Name)
		// Please make sure your host can access this test server,
		// otherwise you will get an error.
		// It is recommended to replace a server at this time
		server.PingTest(nil)
		startDownload := time.Now()
		server.DownloadTest()
		downloadDuration := time.Since(startDownload)

		startUpload := time.Now()
		server.UploadTest()
		uploadDuration := time.Since(startUpload)

		log.Printf("Server: %s; Latency: %s, Download: %f, Upload: %f\n", server.Name, server.Latency, server.DLSpeed, server.ULSpeed)
		log.Printf("Jitter: %s\n", server.Jitter)
		log.Printf("Server ID: %s", server.ID)

		downloadSize := downloadDuration.Seconds() * server.DLSpeed
		uploadSize := uploadDuration.Seconds() * server.ULSpeed
		log.Printf("Download Time: %f, Download Size: %f", downloadDuration.Seconds(), downloadSize)
		log.Printf("Upload Time: %f, Upload Size: %f", uploadDuration.Seconds(), uploadSize)
		result := types.SpeedtestResult{
			Isp:              user.Isp,
			Server:           server.Name,
			Latitude:         server.Lat,
			Longitude:        server.Lon,
			Distance:         server.Distance,
			Latency:          server.Latency.Milliseconds(),
			Jitter:           server.Jitter.Microseconds(),
			Download:         server.DLSpeed,
			Upload:           server.ULSpeed,
			DownloadDuration: downloadDuration.Seconds(),
			DownloadSize:     downloadSize,
			UploadDuration:   uploadDuration.Seconds(),
			UploadSize:       uploadSize,
		}
		service.repository.SaveSpeedtestResults(result)
		server.Context.Reset() // reset counter
	}
}
