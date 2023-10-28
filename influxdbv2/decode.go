package influxdbv2

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/mitchellh/mapstructure"
)

func decode(influxResult *api.QueryTableResult, outData interface{}) error {
	influxData := make([]map[string]interface{}, 0)

	var tags []string

	for influxResult.Next() {
		//fmt.Println("*************************************************************")

		if influxResult.TableChanged() {
			//fmt.Printf("table: %s\n", influxResult.TableMetadata().String())
			for _, col := range influxResult.TableMetadata().Columns() {
				//fmt.Printf("Column: %s\n", col.String())
				if col.IsGroup() {
					if !strings.HasPrefix(col.Name(), "_") {
						tags = append(tags, col.Name())
					}
				}
			}
			fmt.Println("Tags:", tags)
		}

		record := influxResult.Record()
		//fmt.Printf("Time: %v\n", record.Time())
		//fmt.Printf("Values: %v\n", record.Values())
		//fmt.Printf("%v: %v\n", record.Field(), record.Value())

		r := make(map[string]interface{})
		for k, v := range record.Values() {
			if strings.HasPrefix(k, "_") {
				if k != "_field" && k != "_value" {
					r[k] = v
				}
			} else {
				r[k] = v
			}
		}

		find := false
		for _, d := range influxData {
			if d["_time"].(time.Time) != record.Time() {
				continue
			}
			for _, c := range tags {
				if d[c] != record.ValueByKey(c) {
					continue
				}
			}

			d[record.Field()] = record.Value()
			find = true
			break
		}

		if !find {
			r[record.Field()] = record.Value()
			influxData = append(influxData, r)
		}
	}

	if influxResult.Err() != nil {
		fmt.Printf("query parsing error: %s\n", influxResult.Err().Error())
		return influxResult.Err()
	}

	config := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           outData,
		TagName:          "influx",
		WeaklyTypedInput: false,
		ZeroFields:       false,
		DecodeHook: func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
			if t == reflect.TypeOf(time.Time{}) && f == reflect.TypeOf("") {
				return time.Parse(time.RFC3339, data.(string))
			}

			return data, nil
		},
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(influxData)
}
