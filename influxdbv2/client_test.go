package influxdbv2

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type EnvSample struct {
	BasicTag    `influx:",squash"`
	Location    string  `influx:"location,tag"`
	Temperature float64 `influx:"temperature"`
	Humidity    float64 `influx:"humidity"`
	ID          string  `influx:"-"`
}

type Stat struct {
	BasicTag `influx:",squash"`
	Unit     string  `influx:"unit,tag"`
	Avg      float64 `influx:"avg"`
	Max      float64 `influx:"max"`
}

var (
	client *Client
	ctx    context.Context
)

func init() {
	ctx = context.Background()
	opt := &Options{
		Address: "http://localhost:8086",
		Token:   "admintoken123",
		Org:     "primary",
		Bucket:  "rand-buck",
	}
	cli := NewClient(opt)
	client = cli
	_ = cli.CreateBucket(ctx, opt.Org, opt.Bucket)
}

func TestNewInfluxClient(t *testing.T) {
	assert.NotNil(t, client)
}

func TestWriteData(t *testing.T) {
	assert.NotNil(t, client)

	var err error

	env := EnvSample{
		BasicTag: BasicTag{
			Measurement: "env",
			Time:        time.Now(),
		},
		Location:    "Rm 243",
		Temperature: 70.0,
		Humidity:    60.0,
		ID:          "12432as32",
	}
	err = client.BlockWriteData(ctx, env)
	assert.Nil(t, err)

	stat := Stat{
		BasicTag: BasicTag{
			Measurement: "stat",
			Time:        time.Now(),
		},
		Unit: "temperature",
		Avg:  24.5,
		Max:  45.0,
	}
	err = client.BlockWriteData(ctx, stat)
	assert.Nil(t, err)
}

func TestQueryData(t *testing.T) {
	assert.NotNil(t, client)

	var samplesRead []EnvSample

	q := `
from(bucket:"rand-buck")
	|> range(start:-30d)
	|> filter(fn:(r) => r._measurement == "env" and r._field == "temperature" or r._field == "humidity")
	|> yield(name: "_results")
`
	err := client.QueryData(ctx, q, &samplesRead)
	assert.Nil(t, err)
}
