package v2

import (
	"context"
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

type usingValue struct {
	value  string
	retain bool
}

type Client struct {
	opt *Options
	cli influxdb2.Client
}

func NewClient(o *Options) *Client {
	c := &Client{}

	cli := influxdb2.NewClientWithOptions(o.Address, o.Token,
		influxdb2.DefaultOptions().SetBatchSize(20),
	)

	c.cli = cli
	c.opt = o

	return c
}

func (c *Client) Close() {
	if c.cli != nil {
		c.cli.Close()
		c.cli = nil
	}
}

func (c *Client) CreateBucket(ctx context.Context, orgName, bucketName string) error {
	org, err := c.cli.OrganizationsAPI().FindOrganizationByName(ctx, orgName)
	if err != nil {
		fmt.Printf("ERROR. Cannot find organization")
		return nil
	}

	bucketsAPI := c.cli.BucketsAPI()
	_, err = bucketsAPI.CreateBucketWithName(ctx, org, bucketName, domain.RetentionRule{EverySeconds: 3600 * 12})
	if err != nil {
		fmt.Printf("Error. Cannot create bucket")
		return err
	}
	return nil
}

func (c *Client) BlockWriteData(ctx context.Context, data interface{}) error {
	writeAPI := c.cli.WriteAPIBlocking(c.opt.Org, c.opt.Bucket)

	ts, tags, fields, measurement, err := encode(data, nil)
	if err != nil {
		return err
	}

	p := influxdb2.NewPoint(measurement, tags, fields, ts)

	return writeAPI.WritePoint(ctx, p)
}

func (c *Client) WriteData(data interface{}) error {
	writeAPI := c.cli.WriteAPI(c.opt.Org, c.opt.Bucket)
	// Read and log errors
	errorsCh := writeAPI.Errors()
	go func() {
		for err := range errorsCh {
			fmt.Printf("write error: %s\n", err.Error())
		}
	}()

	ts, tags, fields, measurement, err := encode(data, nil)
	if err != nil {
		return err
	}

	p := influxdb2.NewPoint(measurement, tags, fields, ts)

	writeAPI.WritePoint(p)

	writeAPI.Flush()

	return nil
}

func (c *Client) QueryData(ctx context.Context, query string, result interface{}) error {
	queryAPI := c.cli.QueryAPI(c.opt.Org)

	var err error

	var influxResult *api.QueryTableResult
	if influxResult, err = queryAPI.Query(ctx, query); err != nil {
		fmt.Printf("ERROR. Cannot serve qeury result: %v \n", err)
		return err
	}
	defer influxResult.Close()

	if err = decode(influxResult, &result); err != nil {
		return err
	}

	if influxResult.Err() != nil {
		fmt.Printf("query parsing error: %s\n", influxResult.Err().Error())
		return influxResult.Err()
	}

	return nil
}
