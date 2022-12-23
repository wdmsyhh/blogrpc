package util

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	crypto_rand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"net/url"
	"reflect"
	"strings"
	"time"

	"blogrpc/core/extension/bson"

	"github.com/spf13/cast"
	"github.com/wulijun/go-php-serialize/phpserialize"
)

const (
	COOKIE_KEY = "r5jO-PYHr50S6EU88Rt9v70FiEwxXvAC"

	RFC3339     = "2006-01-02T15:04:05Z07:00"
	RFC3339Mili = "2006-01-02T15:04:05.999Z07:00"

	MAX_INT64 int64 = 9223372036854775807
	MAX_INT   int   = 4294967295
)

// the function which will automatically convert value
type converter func(reflect.Value) reflect.Value

type loader func(interface{}, reflect.Value, map[string]interface{}, map[string]string)

type transformer func(interface{}, interface{}, map[string]interface{}, map[string]string)

func genLoader(converters ...converter) loader {
	return func(origin interface{}, targetValue reflect.Value, resetFieldMethods map[string]interface{}, diffFields map[string]string) {
		originValue := reflect.ValueOf(origin)
		if originValue.Type().Kind() == reflect.Ptr {
			originValue = originValue.Elem()
		}

		if targetValue.Type().Kind() == reflect.Ptr {
			targetValue = targetValue.Elem()
		}

		originType := originValue.Type()

		targetType := targetValue.Type()

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
						if targetField.Type != value.Type() {
							for _, convert := range converters {
								value = convert(value)
							}
						}
					}
				}

				if value.Type().ConvertibleTo(targetField.Type) {
					value = value.Convert(targetField.Type)
					targetValue.FieldByName(fieldName).Set(value)
				}
			}
		}
	}
}

// genExistLoader only copy the fields in origin
// attention: if the origin field has a bool type and the value is false, this field will not be copied
func genOverwriteLoader(converters ...converter) loader {
	return func(origin interface{}, targetValue reflect.Value, resetFieldMethods map[string]interface{}, diffFields map[string]string) {
		originValue := reflect.ValueOf(origin)
		if originValue.Type().Kind() == reflect.Ptr {
			originValue = originValue.Elem()
		}

		if targetValue.Type().Kind() == reflect.Ptr {
			targetValue = targetValue.Elem()
		}

		originType := originValue.Type()

		targetType := targetValue.Type()

		for i := 0; i < originType.NumField(); i++ {
			originField := originType.Field(i)
			fieldName := originField.Name
			value := originValue.Field(i)
			// 判断为空则跳过拷贝，bool 类型需要重新判断
			if IsZero(value.Interface()) && reflect.ValueOf(value.Interface()).Kind() != reflect.Bool {
				continue
			}
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
				}
			}
		}
	}
}

// genGetterLoader
//
// fieldsMapper 表示 target 中的字段应该访问 origin 的字段
// getFieldMethods 表示取值方法，形参是 origin 中的字段
func genGetterLoader(converters ...converter) loader {
	return func(origin interface{}, targetValue reflect.Value, getFieldMethods map[string]interface{}, fieldsMapper map[string]string) {
		originValue := reflect.ValueOf(origin)
		if originValue.Type().Kind() == reflect.Ptr {
			originValue = originValue.Elem()
		}

		if targetValue.Type().Kind() == reflect.Ptr {
			targetValue = targetValue.Elem()
		}

		originType := originValue.Type()

		targetType := targetValue.Type()

		// 列举 target 字段
		for i := 0; i < targetType.NumField(); i++ {
			targetField := targetType.Field(i)
			fieldName := targetField.Name
			value := reflect.Value{}
			if originFieldName, ok := fieldsMapper[fieldName]; ok {
				if method, ok := getFieldMethods[fieldName]; ok {
					value = originValue.FieldByName(originFieldName)
					f := reflect.ValueOf(method)
					result := f.Call(
						[]reflect.Value{value},
					)
					value = result[0]
					targetValue.FieldByName(fieldName).Set(value.Convert(targetField.Type))
					continue
				} else {
					panic(fmt.Sprintf("can't find getFieldMethod:[%v]", fieldName))
				}
			}

			if originField, ok := originType.FieldByName(fieldName); ok {
				value = originValue.FieldByName(fieldName)
				if len(converters) > 0 {
					for _, convert := range converters {
						value = convert(value)
					}
				}
				if originField.Type.ConvertibleTo(targetField.Type) {
					value = value.Convert(targetField.Type)
					targetValue.FieldByName(fieldName).Set(value)
				}
			}
		}
	}
}

func genWriteEmptyFieldsGetterLoader(converters ...converter) loader {
	return func(origin interface{}, targetValue reflect.Value, getFieldMethods map[string]interface{}, fieldsMapper map[string]string) {
		originValue := reflect.ValueOf(origin)
		if originValue.Type().Kind() == reflect.Ptr {
			originValue = originValue.Elem()
		}

		if targetValue.Type().Kind() == reflect.Ptr {
			targetValue = targetValue.Elem()
		}

		originType := originValue.Type()

		targetType := targetValue.Type()

		for i := 0; i < targetType.NumField(); i++ {
			targetField := targetType.Field(i)
			fieldName := targetField.Name
			value := reflect.Value{}

			if !targetValue.Field(i).IsZero() {
				continue
			}

			if originFieldName, ok := fieldsMapper[fieldName]; ok {
				if method, ok := getFieldMethods[fieldName]; ok {
					value = originValue.FieldByName(originFieldName)
					f := reflect.ValueOf(method)
					result := f.Call(
						[]reflect.Value{value},
					)
					value = result[0]
					targetValue.FieldByName(fieldName).Set(value.Convert(targetField.Type))
					continue
				} else {
					panic(fmt.Sprintf("can't find getFieldMethod:[%v]", fieldName))
				}
			}

			if _, ok := originType.FieldByName(fieldName); ok {
				value = originValue.FieldByName(fieldName)
				if len(converters) > 0 {
					for _, convert := range converters {
						value = convert(value)
					}
				}
				if value.Type().ConvertibleTo(targetField.Type) {
					value = value.Convert(targetField.Type)
					targetValue.FieldByName(fieldName).Set(value)
				}
			}
		}
	}
}

func genOverWriteGetterLoader(converters ...converter) loader {
	return func(origin interface{}, targetValue reflect.Value, getFieldMethods map[string]interface{}, fieldsMapper map[string]string) {
		originValue := reflect.ValueOf(origin)
		if originValue.Type().Kind() == reflect.Ptr {
			originValue = originValue.Elem()
		}

		if targetValue.Type().Kind() == reflect.Ptr {
			targetValue = targetValue.Elem()
		}

		originType := originValue.Type()

		targetType := targetValue.Type()

		for i := 0; i < targetType.NumField(); i++ {
			targetField := targetType.Field(i)
			fieldName := targetField.Name
			value := reflect.Value{}

			if originValue.FieldByName(fieldName).IsZero() {
				continue
			}

			if originFieldName, ok := fieldsMapper[fieldName]; ok {
				if method, ok := getFieldMethods[fieldName]; ok {
					value = originValue.FieldByName(originFieldName)
					f := reflect.ValueOf(method)
					result := f.Call(
						[]reflect.Value{value},
					)
					value = result[0]
					targetValue.FieldByName(fieldName).Set(value.Convert(targetField.Type))
					continue
				} else {
					panic(fmt.Sprintf("can't find getFieldMethod:[%v]", fieldName))
				}
			}

			if _, ok := originType.FieldByName(fieldName); ok {
				value = originValue.FieldByName(fieldName)
				if len(converters) > 0 {
					for _, convert := range converters {
						value = convert(value)
					}
				}
				if value.Type().ConvertibleTo(targetField.Type) {
					value = value.Convert(targetField.Type)
					targetValue.FieldByName(fieldName).Set(value)
				}
			}
		}
	}
}

func genTransformer(load loader) transformer {
	return func(origin, target interface{}, resetFieldMethods map[string]interface{}, diffFields map[string]string) {
		if reflect.Slice != reflect.TypeOf(origin).Kind() {
			load(origin, reflect.ValueOf(target).Elem(), resetFieldMethods, diffFields)
		} else {
			originValue := reflect.ValueOf(origin)
			length := originValue.Len()
			newSlice := reflect.MakeSlice(reflect.TypeOf(target).Elem(), length, length)
			targetSliceValue := reflect.ValueOf(target).Elem()
			targetSliceValue.Set(newSlice)
			for i := 0; i < length; i++ {
				value := targetSliceValue.Index(i)
				// If the type of the value is a pointer, nil cannot be set directly
				// We set a zero value for it first
				if value.Kind() == reflect.Ptr {
					basicTargetType := value.Type().Elem()
					zeroValue := reflect.New(basicTargetType)
					value.Set(zeroValue)
				}

				load(originValue.Index(i).Interface(), value, resetFieldMethods, diffFields)
			}
		}
	}
}

// Note that type of objectIdhex string will convert to objectId.
func getAutomaticConvertValueObjectId(value reflect.Value) reflect.Value {
	if stringValue, ok := value.Interface().(string); ok {
		if bson.IsObjectIdHex(stringValue) {
			return reflect.ValueOf(bson.ObjectIdHex(stringValue))
		}
	}

	return value
}

// Note that type of time.Time and bson.ObjectId will automatically convert to int64 and string.
func getAutomaticConvertValue(value reflect.Value) reflect.Value {
	if objectIdValue, ok := value.Interface().(bson.ObjectId); ok {
		return reflect.ValueOf(objectIdValue.Hex())
	}

	if timeValue, ok := value.Interface().(time.Time); ok {
		if timeValue.Unix() > 0 {
			return reflect.ValueOf(timeValue.Unix())
		} else {
			return reflect.ValueOf(int64(0))
		}
	}

	return value
}

// similar to getAutomaticConvertValue, but it will convert time to RFC3339
func getAutomaticConvertValueRFC3339(value reflect.Value) reflect.Value {
	if objectIdValue, ok := value.Interface().(bson.ObjectId); ok {
		return reflect.ValueOf(objectIdValue.Hex())
	}

	if timeValue, ok := value.Interface().(time.Time); ok {
		if timeValue.Unix() > 0 {
			return reflect.ValueOf(timeValue.Format(RFC3339))
		} else {
			return reflect.ValueOf("")
		}
	}

	return value
}

var loadProtoBuf = genLoader(getAutomaticConvertValue)
var loadProtoBufRFC3339 = genLoader(getAutomaticConvertValueRFC3339)
var load = genLoader()

// Deprecated: Use copier.Instance.From.CopyTo instead.
var TransformExistFields = genTransformer(loadProtoBufOverwrite)
var loadProtoBufOverwrite = genOverwriteLoader(getAutomaticConvertValue)

// Deprecated: Use copier.Instance.From.CopyTo instead.
var TransformFieldsByGetter = genTransformer(loadProtoBufGetter)

var loadProtoBufGetter = genGetterLoader(getAutomaticConvertValue)

// Deprecated: Use copier.Instance.From.CopyTo instead.
var TransformFieldsModelToModelByGetter = genTransformer(genGetterLoader(getAutomaticConvertValueObjectId))

// Deprecated: Use copier.Instance.From.CopyTo instead.
var TransformExistFieldsByGetter = genTransformer(loadProtoBufWriteEmptyFieldsGetter)
var loadProtoBufWriteEmptyFieldsGetter = genWriteEmptyFieldsGetterLoader(getAutomaticConvertValueRFC3339)

// Deprecated: Use copier.Instance.From.CopyTo instead.
var TransformOverWriteFieldsModelToModelByGetter = genTransformer(genOverWriteGetterLoader(getAutomaticConvertValueObjectId))

// Enhanced CopyFields(origin, target interface{})
// Support to change the field name and reset field value through a custom function.
// resetFieldMethods, the key is the field name, value is reset function.
// Note that type of time.Time and bson.ObjectId will automatically convert.
// diffFields, the key is old field name. value is new filed name.
//
// Deprecated: Use copier.Instance.From.CopyTo instead.
var TransformFields = genTransformer(loadProtoBuf)

// Deprecated: Use copier.Instance.From.CopyTo instead.
var TransformFieldsRFC3339 = genTransformer(loadProtoBufRFC3339)

// Deprecated: Use copier.Instance.From.CopyTo instead.
var TransformFieldsModelToModel = genTransformer(genLoader(getAutomaticConvertValueObjectId, getAutomaticConvertValue))

// Similar to TransformFields, but this function will not automatically
// convert values
//
// Deprecated: Use copier.Instance.From.CopyTo instead.
var TransformFieldsWithoutAutoConvertion = genTransformer(load)

// Copy origin value and stores the result in the value pointed to by target.
// Origin can be slice or struct.
// Target must be a pointer, and not equal to nil.
//
// Deprecated: Use copier.Instance.From.CopyTo instead.
func CopyFields(origin, target interface{}) {
	TransformFields(origin, target, map[string]interface{}{}, map[string]string{})
}

// Deprecated: Use copier.Instance.From.CopyTo instead.
func CopyFieldsRFC3339(origin, target interface{}) {
	TransformFieldsRFC3339(origin, target, map[string]interface{}{}, map[string]string{})
}

// This function is written for struct, all the target and origins are
// supporsed to be struct.
//
// Deprecated: Use copier.Instance.From.CopyTo instead.
func CopyFieldsWithoutConvert(origin, target interface{}) {
	TransformFieldsWithoutAutoConvertion(origin, target, map[string]interface{}{}, map[string]string{})
}

func CopyByJson(origin, target interface{}) error {
	bytes, err := json.Marshal(origin)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, target)
}

func BaseConvert(number string, fromBase, toBase int) (string, bool) {
	if fromBase < 2 || toBase > 36 {
		return "", false
	}

	str := big.Int{}
	_, success := str.SetString(number, fromBase)
	if !success {
		return "", false
	}

	return str.Text(toBase), true
}

func GenUniqueId() string {
	objectId := bson.NewObjectId().Hex()
	id, _ := BaseConvert(objectId, 16, 36)

	// add 5 extra random 36 base digits after id
	// in order to avoid being guessed
	randomNumber := GenRandomNumber(5, 36)

	return fmt.Sprintf("%s%s", id, randomNumber)
}

// due to the "number" can be 2-36 base, so the result can only be a string
func GenRandomNumber(numLength, base int) string {
	maxNumber := int64(math.Pow(float64(base), float64(numLength)))
	nBig, _ := crypto_rand.Int(crypto_rand.Reader, big.NewInt(maxNumber))
	randomNumber := nBig.Int64()

	convertedNumber, _ := BaseConvert(cast.ToString(randomNumber), 10, base)

	// pad the generated number with 0, in case the randomNumber too small
	padFormat := fmt.Sprintf("%s%d%s", "%0", numLength, "s")

	paddedNumber := fmt.Sprintf(padFormat, convertedNumber)

	return paddedNumber
}

func GetCaseInsensitiveStrRegex(str string) bson.RegEx {
	return bson.RegEx{
		Pattern: "^" + FormatRegexStr(str) + "$",
		Options: "i",
	}
}

func FormatRegexStr(str string) string {
	oldnews := []string{
		"\\", "\\\\",
		"*", "\\*",
		".", "\\.",
		"?", "\\?",
		"+", "\\+",
		"$", "\\$",
		"^", "\\^",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"{", "\\{",
		"}", "\\}",
		"|", "\\|",
		"/", "\\/",
	}
	return strings.NewReplacer(oldnews...).Replace(str)
}

func EncryptSha256(pData *[]byte, secretKey string) string {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write(*pData)
	bytes := mac.Sum(nil)
	return hex.EncodeToString(bytes)
}

func ValidateCookieData(data string) string {
	dataTest := []byte{}
	hashTest := EncryptSha256(&dataTest, COOKIE_KEY)
	if len(data) > len(hashTest) {
		hash := data[0:len(hashTest)]
		rawData := data[len(hashTest):]

		rawDataByte := []byte(rawData)
		calculatedHash := EncryptSha256(&rawDataByte, COOKIE_KEY)

		if calculatedHash == hash {
			return rawData
		}
	}

	return ""
}

func DecodeAccessToken(token string) string {
	decodedValue, _ := url.QueryUnescape(token)
	serializedToken := ValidateCookieData(decodedValue)
	if serializedToken != "" {
		unserializedToken, err := phpserialize.Decode(serializedToken)
		if err == nil {
			tokenMap := unserializedToken.(map[interface{}]interface{})
			for k, v := range tokenMap {
				if cast.ToInt(k) == 1 {
					return cast.ToString(v)
				}
			}
		}
	}

	return ""
}

// judge whether the params is zero value
func IsZero(any interface{}) bool {
	v := reflect.ValueOf(any)
	return IsEmpty(v)
}

func IsEmpty(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}

	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

// Avoid cast.ToInt() converts based on 8 when string started with zero. https://github.com/spf13/cast/issues/56
func StringToInt(str string) int {
	return cast.ToInt(strings.TrimLeft(str, "0"))
}

type RetriedFunction func() error

func RetryWrapper(times int, fn func() error) func() error {
	initialCount := 1

	var retryFn func(count int) error
	retryFn = func(count int) error {
		err := fn()
		if err == nil {
			return nil
		}
		if count >= times {
			return err
		}

		count++
		return retryFn(count)
	}

	return func() error {
		return retryFn(initialCount)
	}
}

func StrInArray(search string, items *[]string) bool {
	if items == nil {
		return false
	}
	contains := false
	for _, item := range *items {
		if item == search {
			contains = true
			break
		}
	}
	return contains
}

func GetInt64PointValue(value int64) *int64 {
	return &value
}

func LowercaseFirst(word string) string {
	length := len(word)
	if length == 0 {
		return ""
	}
	remaining := word[1:]
	first := strings.ToLower(string(word[0]))
	return strings.Join([]string{first, remaining}, "")
}

// IntersectStringSlice compare two string slice, return values intersect
func IntersectStringSlice(target, dest []string) []string {
	var intersects []string
	for _, v := range target {
		if StrInArray(v, &dest) {
			intersects = append(intersects, v)
		}
	}

	return intersects
}

// this function will tell you are all fields in element.
//
// the fields support array.
// the element support map,array,struct.
//
// attention: if is struct, the field type is only support string to check if struct has the field.
//
// the return arg[0] is the exist in fields, the args[1] is inexistence in fields.
func GetFieldsInAndNotInElement(fields interface{}, element interface{}) ([]interface{}, []interface{}) {
	slice, len := ConvertToSlice(fields)
	var existFields, inexistenceFields []interface{}
	if len == 0 {
		return nil, slice
	}
	for _, item := range slice {
		switch reflect.TypeOf(item).Kind() {
		// if you want more basic type, you can add here
		case reflect.String, reflect.Int, reflect.Ptr, reflect.Struct, reflect.Int32, reflect.Int64, reflect.Uint64:
			if IsContains(item, element) {
				existFields = append(existFields, item)
			} else {
				inexistenceFields = append(inexistenceFields, item)
			}
		default:
			panic(fmt.Sprintf("not support [%v], need to support.", reflect.TypeOf(item).Kind()))
		}
	}
	return existFields, inexistenceFields
}

func IsContains(item interface{}, element interface{}) bool {
	switch reflect.TypeOf(element).Kind() {
	case reflect.Slice, reflect.Array:
		return IndexOfArray(item, element) != -1
	case reflect.Map:
		return IsInMap(item, element)
	case reflect.Struct:
		// item need to be a string
		return IsStructHasField(item.(string), element)
	default:
		return false
	}
}

func IsStructHasField(fieldName string, element interface{}) bool {
	if reflect.TypeOf(element).Kind() != reflect.Struct {
		panic("not support")
	}
	return reflect.ValueOf(element).FieldByName(fieldName) != reflect.ValueOf(nil)
}

func GetMaxAndMin(nums ...int64) (int64, int64) {
	min, max := int64(math.MaxInt64), int64(math.MinInt64)
	for _, num := range nums {
		if num > max {
			max = num
		}
		if min > num {
			min = num
		}
	}
	return max, min
}

func GetMaxAndMinByRange(start, end int64, nums ...int64) (int64, int64) {
	max, min := GetMaxAndMin(nums...)
	if max > end {
		max = end
	}
	if min < start {
		min = start
	}
	return max, min
}

// DeepCopy copys the in to out deeply. The in is model or slice of models.
// The out is the address of model or the address of slice of models.
func DeepCopy(in, out interface{}) {
	inValue := reflect.ValueOf(in)
	switch inValue.Kind() {
	case reflect.Slice:
		deepCopySlice(in, out)
	case reflect.Ptr:
		inValue = inValue.Elem()
		DeepCopy(inValue.Interface(), out)
	case reflect.Map:
		fallthrough
	case reflect.Struct:
		deepCopyStruct(in, out)
	default:
		outValue := reflect.ValueOf(out).Elem()
		outValue.Set(inValue)
	}
}

func deepCopySlice(in, out interface{}) {
	sliceIn := reflect.ValueOf(in)
	elemtInType := sliceIn.Type().Elem()
	outValue := reflect.ValueOf(out)
	if outValue.Kind() != reflect.Ptr || outValue.Elem().Kind() != reflect.Slice {
		panic("out argument must be a slice address")
	}
	sliceOut := outValue.Elem()
	elemtOutType := sliceOut.Type().Elem()

	if elemtInType != elemtOutType {
		panic("in and out argument must be the same slice element type")
	}
	sliceOut = reflect.MakeSlice(sliceIn.Type(), sliceOut.Len(), sliceIn.Cap())
	for i := 0; i < sliceIn.Len(); i++ {
		elemp := reflect.New(elemtOutType)
		DeepCopy(sliceIn.Index(i).Interface(), elemp.Interface())

		sliceOut = reflect.Append(sliceOut, elemp.Elem())
	}
	outValue.Elem().Set(sliceOut.Slice(0, sliceIn.Len()))
}

func deepCopyStruct(in, out interface{}) {
	inBytes, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(inBytes, out)
	if err != nil {
		panic(err)
	}
}

// 将数组分成指定个子数组
func SeparateSlice(array interface{}, count int) [][]interface{} {
	kind := reflect.TypeOf(array).Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		panic(fmt.Sprintf("Need slice or array type but got [%v]", kind))
	}

	slice, length := ConvertToSlice(array)

	if length < count {
		count = length
	}
	if count == 0 {
		return nil
	}

	blockSize := length / count
	result := make([][]interface{}, count)

	for i := 0; i < count; i++ {
		result[i] = slice[i*blockSize : i*blockSize+blockSize]
	}
	if blockSize*count < length {
		result[count-1] = append(result[count-1], slice[blockSize*count:length]...)
	}
	return result
}

func ZipString(val []byte) (string, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	defer w.Close()
	if _, err := w.Write(val); err != nil {
		return "", err
	}
	if err := w.Flush(); err != nil {
		return "", err
	}
	return b.String(), nil
}

func UnzipString(val string) ([]byte, error) {
	reader, err := gzip.NewReader(strings.NewReader(val))
	if err != nil {
		return []byte{}, err
	}
	defer reader.Close()
	result, err := ioutil.ReadAll(reader)
	if err != nil && !strings.Contains(err.Error(), "unexpected EOF") {
		return []byte{}, err
	}
	return result, nil
}
