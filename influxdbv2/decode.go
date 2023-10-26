package influxdbv2

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/mitchellh/mapstructure"
)

func decode(influxResult *api.QueryTableResult, result interface{}) error {
	influxData := make([]map[string]interface{}, 0)

	for influxResult.Next() {
		//fmt.Println("*************************************************************")

		//position := influxResult.TableMetadata().Position()
		//columns := influxResult.TableMetadata().Columns()
		record := influxResult.Record()

		//fmt.Println("字段：", position, columns)
		//fmt.Printf("Time: %v\n", record.Time())
		//fmt.Printf("ContainerName: %v\n", record.Values())
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
			if d["_time"].(time.Time) == record.Time() {
				d[record.Field()] = record.Value()
				find = true
				break
			}
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
		Result:           result,
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
