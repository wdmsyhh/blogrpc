package copier

import (
	"blogrpc/core/extension/bson"
	"blogrpc/core/util"
	"reflect"
	"time"
)

func RFC3339Convertor(m Mapper) {
	registerRFC3339TimeToStringConverter(m)
	registerStringToRFC3339TimeConverter(m)
}

func ObjectIdConvertor(m Mapper) {
	registerObjectIdToStringConverter(m)
	registerStringToObjectIdConverter(m)
}

func registerRFC3339TimeToStringConverter(m Mapper) {
	m.RegisterConverter(
		Target{
			From: reflect.TypeOf(time.Time{}),
			To:   reflect.TypeOf(""),
		},
		func(from reflect.Value, _ reflect.Type) (reflect.Value, error) {
			if timeValue, ok := from.Interface().(time.Time); ok {
				if timeValue.Unix() > 0 {
					return reflect.ValueOf(timeValue.Format(util.RFC3339)), nil
				} else {
					return reflect.ValueOf(""), nil
				}
			}
			return from, nil
		},
	)
}

func registerStringToRFC3339TimeConverter(m Mapper) {
	m.RegisterConverter(
		Target{
			From: reflect.TypeOf(""),
			To:   reflect.TypeOf(time.Time{}),
		},
		func(from reflect.Value, _ reflect.Type) (reflect.Value, error) {
			if str, ok := from.Interface().(string); ok {
				t, err := time.Parse(util.RFC3339Mili, str)
				if err != nil {
					return reflect.ValueOf(time.Time{}), nil
				}
				return reflect.ValueOf(t), nil
			}
			return from, nil
		},
	)
}

func registerObjectIdToStringConverter(m Mapper) {
	m.RegisterConverter(
		Target{
			From: reflect.TypeOf(bson.ObjectId("")),
			To:   reflect.TypeOf(""),
		},
		func(from reflect.Value, _ reflect.Type) (reflect.Value, error) {
			if value, ok := from.Interface().(bson.ObjectId); ok {
				return reflect.ValueOf(value.Hex()), nil
			}
			return from, nil
		},
	)
}

func registerStringToObjectIdConverter(m Mapper) {
	m.RegisterConverter(
		Target{
			From: reflect.TypeOf(""),
			To:   reflect.TypeOf(bson.ObjectId("")),
		},
		func(from reflect.Value, _ reflect.Type) (reflect.Value, error) {
			if value, ok := from.Interface().(string); ok && bson.IsObjectIdHex(value) {
				return reflect.ValueOf(bson.ObjectIdHex(value)), nil
			}
			return reflect.ValueOf(bson.ObjectId("")), nil
		},
	)
}
