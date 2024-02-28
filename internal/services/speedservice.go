package services

import (
	"log"
	"speedstar/internal/db"
	"speedstar/internal/types"

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
	user, _ := service.client.FetchUserInfo()
	log.Printf("ISP %s\n", user.Isp)

	serverList, _ := service.client.FetchServers()
	targets, _ := serverList.FindServer([]int{})
	for _, server := range targets {
		log.Printf("Server: %s", server.Name)
		// Please make sure your host can access this test server,
		// otherwise you will get an error.
		// It is recommended to replace a server at this time
		server.PingTest(nil)
		server.DownloadTest()
		server.UploadTest()
		log.Printf("Server: %s; Latency: %s, Download: %f, Upload: %f\n", server.Name, server.Latency, server.DLSpeed, server.ULSpeed)
		log.Printf("Jitter: %s\n", server.Jitter)

		result := types.SpeedtestResult{
			Isp:       user.Isp,
			Server:    server.Name,
			Latitude:  server.Lat,
			Longitude: server.Lon,
			Distance:  server.Distance,
			Latency:   server.Latency.Milliseconds(),
			Jitter:    server.Jitter.Milliseconds(),
			Download:  server.DLSpeed,
			Upload:    server.ULSpeed,
		}
		service.repository.SaveSpeedtestResults(result)
		server.Context.Reset() // reset counter
	}
}
