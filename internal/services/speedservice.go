package services

import (
	"fmt"
	"log"

	"github.com/showwin/speedtest-go/speedtest"
)

type ISpeedTestService interface {
	RunSpeedtest()
}

type SpeedTestService struct {
	client *speedtest.Speedtest
}

func NewSpeedtestService() ISpeedTestService {
	var speedtestClient = speedtest.New()
	return &SpeedTestService{
		client: speedtestClient,
	}
}

func (service SpeedTestService) RunSpeedtest() {
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
		server.Context.Reset() // reset counter
	}
}

func TestSpeed() {
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
