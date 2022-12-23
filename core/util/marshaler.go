package util

import (
	"blogrpc/core/extension/bson"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"reflect"
)

var (
	jsonpbMarshaler = &jsonpb.Marshaler{
		EnumsAsInts:  true,
		EmitDefaults: true,
		OrigName:     true,
	}
)

func MarshalPB2JSON(x interface{}) string {
	if x == nil || reflect.ValueOf(x).IsNil() {
		return fmt.Sprintf("<nil>")
	}

	pb, ok := x.(proto.Message)
	if !ok {
		return fmt.Sprintf("Marshal to json error: not a proto message")
	}

	var buf bytes.Buffer
	if err := jsonpbMarshaler.Marshal(&buf, pb); err != nil {
		return fmt.Sprintf("Marshal to json error: %s", err.Error())
	}
	return buf.String()
}

func MarshalModel2BSON(x interface{}) (bson.M, error) {
	data, err := bson.Marshal(x)
	if err != nil {
		return nil, err
	}

	bsonMap := bson.M{}
	err = bson.Unmarshal(data, &bsonMap)
	if err != nil {
		return nil, err
	}

	return bsonMap, nil
}

func MarshalInterfaceToString(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		return ""
	}

	return string(b)
}

func UnmarshalJsonString(s string) map[string]interface{} {
	m := map[string]interface{}{}
	b := []byte(s)
	json.Unmarshal(b, &m)
	return m
}
