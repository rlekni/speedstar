package db

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"speedstar/internal/types"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var (
	dbUrl    = os.Getenv("INFLUXDB_URL")
	dbToken  = os.Getenv("INFLUXDB_TOKEN")
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
		AddTag("key", "value").
		AddField("isp", result.Isp).
		AddField("server", result.Server).
		AddField("latitude", result.Latitude).
		AddField("longitude", result.Longitude).
		AddField("distance", result.Distance).
		AddField("latency", result.Latency).
		AddField("jitter", result.Jitter).
		AddField("download", result.Download).
		AddField("upload", result.Upload).
		SetTime(time.Now())

	// write asynchronously
	writeAPI.WritePoint(point)

	// Force all unwritten data to be sent
	writeAPI.Flush()
}

func InfluxdbWrites() {
	fmt.Printf("INFLUXDB_ORG: %s\n", dbOrg)
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

func InfluxdbWritesAsync() {
	fmt.Printf("INFLUXDB_ORG: %s\n", dbOrg)
	// Create a new client using an InfluxDB server base URL and an authentication token
	client := influxdb2.NewClientWithOptions(dbUrl, dbToken, influxdb2.DefaultOptions().SetBatchSize(20))
	// Ensures background processes finishes
	defer client.Close()
	// Get non-blocking write client
	writeAPI := client.WriteAPI(dbOrg, dbBucket)
	// write some points
	for i := 0; i < 100; i++ {
		// create point
		p := influxdb2.NewPoint(
			"system",
			map[string]string{
				"id":       fmt.Sprintf("rack_%v", i%10),
				"vendor":   "AWS",
				"hostname": fmt.Sprintf("host_%v", i%100),
			},
			map[string]interface{}{
				"temperature": rand.Float64() * 80.0,
				"disk_free":   rand.Float64() * 1000.0,
				"disk_total":  (i/10 + 1) * 1000000,
				"mem_total":   (i/100 + 1) * 10000000,
				"mem_free":    rand.Uint64(),
			},
			time.Now())
		// write asynchronously
		writeAPI.WritePoint(p)
	}

	// Force all unwritten data to be sent
	writeAPI.Flush()
}
