package util

import (
	"reflect"
	"strings"
	"time"
)

type formatLoader func(reflect.Value, reflect.Value, map[string]interface{}, map[string]string)

func genFormatLoader(converters ...converter) formatLoader {
	var newLoader formatLoader

	newLoader = func(originValue, targetValue reflect.Value, resetFieldMethods map[string]interface{}, diffFields map[string]string) {
		if originValue.Kind() == reflect.Ptr {
			if originValue.IsNil() {
				return
			}
			originValue = originValue.Elem()
		}

		if targetValue.Type().Kind() == reflect.Ptr {
			if targetValue.IsNil() {
				targetValue.Set(reflect.New(targetValue.Type().Elem()))
			}
			targetValue = targetValue.Elem()
		}

		originType := originValue.Type()
		targetType := targetValue.Type()

		switch originType.Kind() {
		case reflect.Struct:

			for i := 0; i < originType.NumField(); i++ {
				originField := originType.Field(i)
				fieldName := originField.Name
				value := originValue.Field(i)
				if method, ok := resetFieldMethods[fieldName]; ok {
					f := reflect.ValueOf(method)
					result := f.Call(
						[]reflect.Value{value},
					)
					value = result[0]
				}

				if newFiledName, ok := diffFields[fieldName]; ok {
					fieldName = newFiledName
				}

				if targetField, ok := targetType.FieldByName(fieldName); ok {
					if len(converters) > 0 {
						_, hasResetFieldMethod := resetFieldMethods[fieldName]

						if !hasResetFieldMethod {
							for _, convert := range converters {
								value = convert(value)
							}
						}
					}

					if value.Type().ConvertibleTo(targetField.Type) {
						value = value.Convert(targetField.Type)
						targetValue.FieldByName(fieldName).Set(value)
						continue
					}

					deepResetFieldMethods := goDeeperReset(resetFieldMethods, fieldName)
					deepDiffFields := goDeeperRename(diffFields, fieldName)
					newLoader(value, targetValue.FieldByName(fieldName), deepResetFieldMethods, deepDiffFields)
				}
			}
		case reflect.Slice:
			length := originValue.Len()
			newSlice := reflect.MakeSlice(targetType, length, length)
			targetValue.Set(newSlice)
			for i := 0; i < length; i++ {
				value := targetValue.Index(i)
				// If the type of the value is a pointer, nil cannot be set directly
				// We set a zero value for it first
				if value.Kind() == reflect.Ptr {
					basicTargetType := value.Type().Elem()
					zeroValue := reflect.New(basicTargetType)
					value.Set(zeroValue)
				}

				newLoader(originValue.Index(i), value, resetFieldMethods, diffFields)
			}
		default:
			if len(converters) > 0 {
				for _, convert := range converters {
					originValue = convert(originValue)
				}
			}

			if originType.ConvertibleTo(targetType) {
				targetValue.Set(originValue)
			}
		}

	}

	return newLoader
}

func goDeeperReset(m map[string]interface{}, field string) map[string]interface{} {
	result := map[string]interface{}{}
	prefix := field + "."
	for key, value := range m {

		if !strings.HasPrefix(key, prefix) {
			continue
		}

		key := strings.TrimPrefix(key, prefix)
		result[key] = value
	}

	return result
}

func goDeeperRename(m map[string]string, field string) map[string]string {
	result := map[string]string{}
	prefix := field + "."
	for key, value := range m {

		if !strings.HasPrefix(key, prefix) {
			continue
		}

		key := strings.TrimPrefix(key, prefix)
		result[key] = value
	}

	return result
}

func genFormatTransformer(load formatLoader) transformer {
	return func(origin, target interface{}, resetFieldMethods map[string]interface{}, diffFields map[string]string) {
		originValue := reflect.ValueOf(origin)
		targetValue := reflect.ValueOf(target)
		if originValue.Type().Kind() == reflect.Ptr {
			originValue = originValue.Elem()
		}
		if targetValue.Type().Kind() == reflect.Ptr {
			targetValue = targetValue.Elem()
		}

		load(originValue, targetValue, resetFieldMethods, diffFields)
	}
}

var formatLoaderRFC3339 = genFormatLoader(getAutomaticConvertValueRFC3339)

var FormatRFC3339 = genFormatTransformer(formatLoaderRFC3339)

var formatLoaderIntDate = genFormatLoader(getAutomaticConvertValue)
var FormatIntDate = genFormatTransformer(formatLoaderIntDate)

var clone = genFormatLoader()
var CloneTransform = genFormatTransformer(clone)

var format = genFormatLoader(getAutomaticConvertValue)
var Format = genFormatTransformer(format)

func CopyRFC3339(origin, target interface{}) {
	FormatRFC3339(origin, target, map[string]interface{}{}, map[string]string{})
}

func Clone(origin, target interface{}) {
	CloneTransform(origin, target, map[string]interface{}{}, map[string]string{})
}

func ParseRFC3339(t string) time.Time {
	paresdTime, _ := time.Parse(RFC3339Mili, t)
	return paresdTime
}
