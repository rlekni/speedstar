package types

type SpeedtestResult struct {
	Isp       string
	Server    string
	Latitude  string
	Longitude string
	Distance  float64
	Latency   int64
	Jitter    int64
	Download  float64
	Upload    float64
}
