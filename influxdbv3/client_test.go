package influxdbv3

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	client *Client
	ctx    context.Context
)

func init() {
	ctx = context.Background()
	opt := &Options{
		Address:  "http://localhost:8086",
		Token:    "admintoken123",
		DataBase: "rand-buck",
	}
	cli := NewClient(opt)
	client = cli
}

func TestQueryData(t *testing.T) {
	assert.NotNil(t, client)

	//var samplesRead []EnvSample

	query := `
        SELECT *
        FROM "stat"
        WHERE
        time >= now() - interval '5 minute'
        AND
        "unit" IN ('temperature')
`
	err := client.QueryData(ctx, query)
	assert.Nil(t, err)
}
