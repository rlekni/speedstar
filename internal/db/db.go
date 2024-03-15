package db

import (
	"os"
	"speedstar/internal/types"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var (
	dbOrg    = os.Getenv("INFLUXDB_ORG")
	dbBucket = os.Getenv("INFLUXDB_BUCKET")
)

type ISpeedtestRepository interface {
	SaveSpeedtestResults(types.SpeedtestResult)
}

type SpeedtestRepository struct {
	client influxdb2.Client
}

func NewSpeedtestRepository(client influxdb2.Client) ISpeedtestRepository {
	return &SpeedtestRepository{
		client: client,
	}
}

func (repo SpeedtestRepository) SaveSpeedtestResults(result types.SpeedtestResult) {
	// Get non-blocking write client
	writeAPI := repo.client.WriteAPI(dbOrg, dbBucket)

	point := influxdb2.NewPointWithMeasurement("speed").
		AddTag("server", result.Server).
		AddField("distance", result.Distance).
		AddField("latency", result.Latency).
		AddField("jitter", result.Jitter).
		AddField("download", result.Download).
		AddField("upload", result.Upload).
		AddField("download_duration", result.DownloadDuration).
		AddField("download_size", result.DownloadSize).
		AddField("upload_duration", result.UploadDuration).
		AddField("upload_size", result.UploadSize).
		SetTime(time.Now())

	// write asynchronously
	writeAPI.WritePoint(point)

	// Force all unwritten data to be sent
	writeAPI.Flush()
}
