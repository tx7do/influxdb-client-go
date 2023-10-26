package influxdbv3

import (
	"context"
	"fmt"
	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
)

type Client struct {
	opt *Options
	cli *influxdb3.Client
}

func NewClient(opt *Options) *Client {
	cli := &Client{
		opt: opt,
	}

	client, _ := influxdb3.New(influxdb3.ClientConfig{
		Host:     opt.Address,
		Token:    opt.Token,
		Database: opt.DataBase,
	})

	cli.cli = client

	return cli
}

func (c *Client) Close() {
	if c.cli != nil {
		_ = c.cli.Close()
		c.cli = nil
	}
}

func (c *Client) QueryData(ctx context.Context, query string) error {
	iterator, err := c.cli.Query(ctx, query)
	if err != nil {
		return err
	}
	for iterator.Next() {
		value := iterator.Value()

		fmt.Printf("avg is %f\n", value["avg"])
		fmt.Printf("max is %f\n", value["max"])
	}
	return nil
}
