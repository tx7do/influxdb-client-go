package influxdbv2

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

func encode(d interface{}, timeField *usingValue) (ts time.Time, tags map[string]string, fields map[string]interface{}, measurement string, err error) {
	tags = make(map[string]string)
	fields = make(map[string]interface{})
	dValue := reflect.ValueOf(d)

	if dValue.Kind() == reflect.Ptr {
		dValue = reflect.Indirect(dValue)
	}

	if dValue.Kind() != reflect.Struct {
		err = errors.New("data must be a struct")
		return
	}

	for i := 0; i < dValue.NumField(); i++ {
		f := dValue.Field(i)
		currField := dValue.Type().Field(i)
		structFieldName := currField.Name

		//fmt.Println("+++++++++++", currField.Type)

		if currField.Type.Kind() == reflect.Struct {
			for j := 0; j < currField.Type.NumField(); j++ {
				sf := currField.Type.Field(j)
				substructFieldName := sf.Name

				parseField(substructFieldName, sf, f.Field(j), timeField,
					&ts, &tags, &fields, &measurement)
			}
			continue
		}

		parseField(structFieldName, currField, f, timeField,
			&ts, &tags, &fields, &measurement)
	}

	if measurement == "" {
		measurement = dValue.Type().Name()
	}

	return
}

func parseField(structFieldName string, currField reflect.StructField, currValue reflect.Value, timeField *usingValue,
	ts *time.Time, tags *map[string]string, fields *map[string]interface{}, measurement *string,
) bool {
	if timeField == nil {
		timeField = &usingValue{"_time", false}
	}

	fieldTag := currField.Tag.Get("influx")
	fieldData := getInfluxFieldTagData(structFieldName, fieldTag)
	//fmt.Println(fieldTag)

	if fieldData.fieldName == "-" {
		return true
	}

	if fieldData.fieldName == timeField.value {
		// TODO error checking
		*ts = currValue.Interface().(time.Time)
		return true
	}

	if fieldData.fieldName == "_measurement" {
		*measurement = currValue.String()
		return true
	}

	if fieldData.isTag {
		(*tags)[fieldData.fieldName] = fmt.Sprintf("%v", currValue)
	}

	if fieldData.isField {
		(*fields)[fieldData.fieldName] = currValue.Interface()
	}

	return true
}
