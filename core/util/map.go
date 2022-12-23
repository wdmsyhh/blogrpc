package util

import (
	"blogrpc/core/extension/bson"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

// This function is written to generate a map.
// key is the field name you want to be the key, array is the map value.
// the array argument support map array or struct array.
func MakeMapper(key string, array interface{}) map[interface{}]interface{} {
	slice, length := ConvertToSlice(array)

	mapper := make(map[interface{}]interface{}, length)

	for _, item := range slice {
		value := GetValueByFiledName(item, key)
		switch reflect.TypeOf(value).Kind() {
		case reflect.Slice, reflect.Array:
			tmpSlice, _ := ConvertToSlice(value)
			mapper[tmpSlice[0]] = item
		default:
			mapper[value] = item
		}
	}
	return mapper
}

// map 的 string 必须是可哈希的，所以如果需要新类型，需要添加到 K 中
func MakeMapperV2[K string | bson.ObjectId | uint64, V any](key string, keyValueInstance K, array []V) map[K]V {
	mapper := make(map[K]V, len(array))
	for _, item := range array {
		value := GetValueByFiledNameV2(item, key, keyValueInstance)
		switch reflect.TypeOf(value).Kind() {
		case reflect.Slice, reflect.Array:
			tmpSlice, _ := ConvertToSlice(value)
			mapper[tmpSlice[0].(K)] = item
		default:
			mapper[value] = item
		}
	}
	return mapper
}

// make mapper from the value of keyField to the value of valueField.
func MakeFieldToFieldMapper(keyField, valueField string, array interface{}) map[interface{}]interface{} {
	slice, length := ConvertToSlice(array)
	mapper := make(map[interface{}]interface{}, length)
	for _, item := range slice {
		key := GetValueByFiledName(item, keyField)
		if key == nil {
			continue
		}
		mapper[key] = GetValueByFiledName(item, valueField)
	}
	return mapper
}

// 将结构体中存在的字段转为 fieldName - fieldValue 形式的 map
//
//	needLowercaseFirst：	fieldName 首字母是否需要小写
//	isBoolFalseAsEmpty：	当 bool 类型的字段为 false 时是否认为是空值
//	hasId：			是否包含 Id：true 则会将 fieldName 转为 _id
type ConvertStructExistBasicFieldsToMapOption struct {
	NeedLowercaseFirst    bool
	IsBoolFalseAsEmpty    bool
	IsIntZeroNeedSet      bool
	HasId                 bool
	IgnoreCreatedAt       bool
	NeedConvertFieldNames []string // 结构体字段名，需大写（如果导出的话）
	AddUpdatedAt          bool
}

func ConvertStructExistBasicFieldsToMap(element interface{}, option ConvertStructExistBasicFieldsToMapOption) map[string]interface{} {
	var elementValue reflect.Value
	if reflect.ValueOf(element).Kind() == reflect.Ptr {
		elementValue = reflect.ValueOf(element).Elem()
	} else {
		elementValue = reflect.ValueOf(element)
	}
	if elementValue.Kind() != reflect.Struct {
		panic("not support" + elementValue.Kind().String())
	}

	result := make(map[string]interface{}, elementValue.Type().NumField())
	elementType := elementValue.Type()
	for i := 0; i < elementType.NumField(); i++ {
		var fieldName string
		if len(option.NeedConvertFieldNames) != 0 && IndexOfArray(elementType.Field(i).Name, option.NeedConvertFieldNames) == -1 {
			continue
		}
		value := elementValue.Field(i).Interface()
		if option.NeedLowercaseFirst {
			fieldName = LowercaseFirst(elementType.Field(i).Name)
		} else {
			fieldName = elementType.Field(i).Name
		}
		if option.HasId && (elementType.Field(i).Name == "Id" || elementType.Field(i).Name == "id") {
			fieldName = "_id"
		}
		if option.IgnoreCreatedAt && (elementType.Field(i).Name == "CreatedAt" || elementType.Field(i).Name == "createdAt") {
			continue
		}

		if IsZero(value) {
			// 当 value 是 bool 类型并且值为 false 时，通过参数 isBoolFalseAsEmpty 判断
			// isBoolFalseAsEmpty 是否 bool 类型 false 算空
			switch reflect.ValueOf(value).Kind() {
			case reflect.Bool:
				if option.IsBoolFalseAsEmpty {
					continue
				}
			case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint32, reflect.Uint64, reflect.Int32:
				if !option.IsIntZeroNeedSet {
					continue
				}
			default:
				continue
			}
		}
		result[fieldName] = value
	}
	if option.AddUpdatedAt {
		result["updatedAt"] = time.Now()
	}
	return result
}

// 读取 element 中 name 字段的值
//
// element 支持结构体、Map
// name 字段名
// ATTENTION：如果获取的字段值为空时不能仅仅判断是否为 Nil，此时 return 的是值为 Nil，类型为原字段类型的 interface 变量。
func GetValueByFiledName(element interface{}, name string) interface{} {
	value := reflect.ValueOf(element)
	var keys []string
	if strings.ContainsAny(name, ".") {
		keys = strings.Split(name, ".")
		firstKey := keys[0]
		keys = append(keys[1:])
		name = strings.Join(keys, ".")
		return GetValueByFiledName(GetValueByFiledName(element, firstKey), name)
	} else {
		switch value.Kind() {
		case reflect.Ptr:
			return value.Elem().FieldByName(name).Interface()
		case reflect.Struct:
			return value.FieldByName(name).Interface()
		case reflect.Map:
			return value.MapIndex(reflect.ValueOf(name)).Interface()
		case reflect.Slice, reflect.Array:
			if value.IsZero() {
				return nil
			} else {
				return ExtractArrayField(name, value.Interface())
			}

		default:
			panic(fmt.Sprintf("not support get value from type[%s]", value.Kind()))
		}
	}
}

func GetValueByFiledNameV2[T any](element any, name string, valueTypeInstance T) T {
	value := reflect.ValueOf(element)
	var keys []string
	if strings.ContainsAny(name, ".") {
		keys = strings.Split(name, ".")
		firstKey := keys[0]
		keys = append(keys[1:])
		name = strings.Join(keys, ".")
		return GetValueByFiledNameV2(GetValueByFiledNameV2(element, firstKey, valueTypeInstance), name, valueTypeInstance)
	} else {
		switch value.Kind() {
		case reflect.Ptr:
			return value.Elem().FieldByName(name).Interface().(T)
		case reflect.Struct:
			return value.FieldByName(name).Interface().(T)
		case reflect.Map:
			return value.MapIndex(reflect.ValueOf(name)).Interface().(T)
		default:
			panic(fmt.Sprintf("not support get value from type[%s]", value.Kind()))
		}
	}
}

// the item type also need the same type with the mapper key
func IsInMap(item interface{}, mapper interface{}) bool {
	if mapper == nil {
		panic("can't be nil")
	}
	if reflect.TypeOf(mapper).Kind() != reflect.Map {
		panic("no supported")
	}
	itemType := reflect.TypeOf(item).Kind()
	iterator := reflect.ValueOf(mapper).MapRange()
	for iterator.Next() {
		if iterator.Key().Interface() == item {
			return true
		} else if reflect.TypeOf(iterator.Key().Interface()).Kind() != itemType {
			panic(fmt.Sprintf("item type [%v] is not match key type [%v]", itemType, reflect.TypeOf(iterator.Key().Interface()).Kind()))
		}
	}
	return false
}

func MakeArrayMapper(key string, array interface{}) map[interface{}][]interface{} {
	slice, length := ConvertToSlice(array)
	mapper := make(map[interface{}][]interface{}, length)
	for _, item := range slice {
		value := GetValueByFiledName(item, key)
		if value == nil {
			continue
		}
		mapper[value] = append(mapper[value], item)
	}
	return mapper
}

func MakeMapperByJudge(key string, array interface{}, opt interface{}) map[interface{}]interface{} {
	f := reflect.ValueOf(opt)
	if f.Kind() != reflect.Func || f.Type().Out(0).Kind() != reflect.Bool {
		panic("arg opt must be func and must return bool.")
	}
	slice, length := ConvertToSlice(array)

	mapper := make(map[interface{}]interface{}, length)
	for _, item := range slice {
		value := GetValueByFiledName(item, key)
		needAdd := f.Call(
			[]reflect.Value{reflect.ValueOf(item)},
		)[0].Interface().(bool)
		if needAdd {
			switch reflect.TypeOf(value).Kind() {
			case reflect.Slice, reflect.Array:
				tmpSlice, _ := ConvertToSlice(value)
				mapper[tmpSlice[0]] = item
			default:
				mapper[value] = item
			}
		}
	}
	return mapper
}

func DecodeMapToStruct(m map[string]interface{}, output interface{}) error {
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   output,
		TagName:  "json",
	}
	decoder, _ := mapstructure.NewDecoder(config)
	return decoder.Decode(m)
}
