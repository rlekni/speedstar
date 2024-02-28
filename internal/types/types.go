package types

type SpeedtestResult struct {
	Isp      string
	Server   string
	Latency  int64
	Download float64
	Upload   float64
}
