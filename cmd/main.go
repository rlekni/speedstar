package main

import (
	"fmt"

	"github.com/showwin/speedtest-go/speedtest"
)

func main() {
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
