package util

import (
	"blogrpc/core/extension/bson"
	"blogrpc/core/util/algorithm"
	"fmt"
	"reflect"
	"sort"

	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// this function help you extract all value of array item named fieldName into a new array.
//
// 该方法可以将数组中指定字段的值提取出来，组成新的切片
func ExtractArrayField(fieldName string, array interface{}) []interface{} {
	slice, len := ConvertToSlice(array)

	fieldValues := make([]interface{}, len)

	for i, item := range slice {
		fieldValues[i] = GetValueByFiledName(item, fieldName)
	}
	return fieldValues
}

func ExtractArrayFieldV2[T any](fieldName string, arrayItemTypeInstance T, array interface{}) []T {
	slice, len := ConvertToSlice(array)

	fieldValues := make([]T, len)

	for i, item := range slice {
		fieldValues[i] = GetValueByFiledNameV2(item, fieldName, arrayItemTypeInstance)
	}
	return fieldValues
}

// 该方法可以将数组中指定字段的值提取出来，有选择的组成新的切片
func ExtractArrayFieldWithJudge(fieldName string, array interface{}, opt interface{}) []interface{} {
	f := reflect.ValueOf(opt)
	if f.Kind() != reflect.Func || f.Type().Out(0).Kind() != reflect.Bool {
		panic("arg opt must be func and must return bool.")
	}

	slice, _ := ConvertToSlice(array)
	var fieldValues []interface{}

	for _, item := range slice {
		needAdd := f.Call(
			[]reflect.Value{reflect.ValueOf(item)},
		)[0].Interface().(bool)
		if needAdd {
			fieldValues = append(fieldValues, GetValueByFiledName(item, fieldName))
		}
	}
	return fieldValues
}

func ExtractArrayStringField(fieldName string, array interface{}) []string {
	slice, len := ConvertToSlice(array)

	fieldValues := make([]string, len)

	for i, item := range slice {
		fieldValues[i] = ToString(GetValueByFiledName(item, fieldName))
	}
	return fieldValues
}

// 将 string 、 ObjectId 和 uint64 类型的空接口切片转为 string 切片
//
// 如果需要更多类型支持，请添加。不在支持列表中的会导致 panic
func ToStringArray(array any) []string {
	slice, len := ConvertToSlice(array)
	result := make([]string, len)
	for i, item := range slice {
		switch item.(type) {
		case string:
			result[i] = item.(string)
		case bson.ObjectId:
			result[i] = item.(bson.ObjectId).Hex()
		case primitive.ObjectID:
			result[i] = item.(primitive.ObjectID).Hex()
		case uint64, float64, float32, int:
			result[i] = cast.ToString(item)
		default:
			panic(fmt.Sprintf("no support type [%v] convert to string", reflect.TypeOf(item).Kind()))
		}
	}
	return result
}

// 将 string 和 ObjectId 类型的空接口切片转为 ObjectId 切片
func ToObjectIdArray(array interface{}) []bson.ObjectId {
	slice, _ := ConvertToSlice(array)
	result := []bson.ObjectId{}
	var id bson.ObjectId
	for _, item := range slice {
		id = ToObjectId(item)
		if id.Hex() != "" {
			result = append(result, id)
		}
	}
	return result
}

// 如果需要更多类型支持，请添加。不在支持列表中的会导致 panic
func ToObjectId(value interface{}) bson.ObjectId {
	var id bson.ObjectId
	switch value.(type) {
	case bson.ObjectId:
		id = value.(bson.ObjectId)
	case primitive.ObjectID:
		id = bson.ObjectIdHex(value.(primitive.ObjectID).Hex())
	case string:
		id = bson.ObjectIdHex(value.(string))
	default:
		panic(fmt.Sprintf("no support type [%v] convert to bson.ObjectId", reflect.TypeOf(value).Kind()))
	}
	return id
}

// 将 string 和 ObjectId 类型的空接口切片转为 ObjectId 切片
//
// 如果需要更多类型支持，请添加。不在支持列表中的会导致 panic
func ConvertToObjectIdArray(array []interface{}) []bson.ObjectId {
	result := []bson.ObjectId{}
	var id bson.ObjectId
	for _, item := range array {
		switch item.(type) {
		case bson.ObjectId:
			id = item.(bson.ObjectId)
		case string:
			id = bson.ObjectIdHex(item.(string))
		default:
			panic(fmt.Sprintf("no support type [%v] convert to bson.ObjectId", reflect.TypeOf(item).Kind()))
		}
		if id.Hex() != "" {
			result = append(result, id)
		}
	}
	return result
}

func ToInt32Array(array interface{}) []int32 {
	slice, len := ConvertToSlice(array)
	result := make([]int32, len)
	for i, item := range slice {
		switch item.(type) {
		case int32:
			result[i] = item.(int32)
		default:
			panic(fmt.Sprintf("no support type [%v] convert to int32", reflect.TypeOf(item).Kind()))
		}
	}
	return result
}

func ToUint64Array(array interface{}) []uint64 {
	slice, len := ConvertToSlice(array)
	result := make([]uint64, len)
	for i, item := range slice {
		switch item.(type) {
		case uint64:
			result[i] = item.(uint64)
		case int64:
			result[i] = uint64(item.(int64))
		default:
			panic(fmt.Sprintf("no support type [%v] convert to uint64", reflect.TypeOf(item).Kind()))
		}
	}
	return result
}

func ConvertToInt32Array(array *[]interface{}) []int32 {
	result := make([]int32, len(*array))
	for i, item := range *array {
		switch item.(type) {
		case int32:
			result[i] = item.(int32)
		default:
			panic(fmt.Sprintf("no support type [%v] convert to int32", reflect.TypeOf(item).Kind()))
		}
	}
	return result
}

// 读取多个切片的交集，类型需要一致
// 2020: BenchmarkGetArraysIntersection-8    	    1576	    763640 ns/op	  431162 B/op	   10208 allocs/op
// 2022: BenchmarkGetArraysIntersection-8   	   43219	     34338 ns/op	   23488 B/op	     427 allocs/op
func GetArraysIntersection(arrays ...interface{}) []interface{} {
	if len(arrays) == 0 {
		return nil
	}
	set := algorithm.CSet.InstanceFromStructSlice(arrays[0])
	for i := range arrays {
		if i == 0 {
			continue
		}
		result := algorithm.CSet.Instance()
		tmpArray, _ := ConvertToSlice(arrays[i])
		for j := range tmpArray {
			if set.Has(tmpArray[j]) {
				result.Insert(tmpArray[j])
			}
		}
		if result.Size() == 0 {
			return nil
		}
		set = result
	}
	return set.ToArray()
}

// 读取多个切片的交集，类型需要一致，不再支持结构体数组
// 2020: BenchmarkGetArraysIntersection-8    	    1576	    763640 ns/op	  431162 B/op	   10208 allocs/op
// 2022: BenchmarkGetArraysIntersection-8   	   43219	     34338 ns/op	   23488 B/op	     427 allocs/op
func GetArraysIntersectionV2[T algorithm.SetType](arrays ...[]T) []T {
	if len(arrays) == 0 {
		return nil
	}
	set := (&algorithm.SetV2[T]{}).InstanceFromSlice(&arrays[0])
	for i := range arrays {
		if i == 0 {
			continue
		}
		result := (&algorithm.SetV2[T]{}).Instance()
		tmpArray, _ := ConvertToSlice(arrays[i])
		for j := range tmpArray {
			if set.Has(tmpArray[j].(T)) {
				result.Insert(tmpArray[j].(T))
			}
		}
		if result.Size() == 0 {
			return nil
		}
		set = result
	}
	return set.ToArray()
}

func ContainsInt64(array *[]int64, element int64) bool {
	return IndexOfArray(element, array) >= 0
}

// 元素在数组中的位置
func IndexOfArray(item interface{}, array interface{}) int {
	slice, _ := ConvertToSlice(array)
	itemType := reflect.TypeOf(item).Kind()
	for i, v := range slice {
		if reflect.TypeOf(v).Kind() != itemType {
			panic(fmt.Sprintf("item type [%v] not equal array item type [%v]! ", itemType, reflect.TypeOf(v).Kind()))
		}
		if v == item {
			return i
		}
	}
	return -1
}

// 元素在 interface 数组中的位置
func IndexOfInterfaceArray(item interface{}, array []interface{}) int {
	itemType := reflect.TypeOf(item).Kind()
	for i, v := range array {
		if reflect.TypeOf(v).Kind() != itemType {
			panic(fmt.Sprintf("item type [%v] not equal array item type [%v]! ", itemType, reflect.TypeOf(v).Kind()))
		}
		if v == item {
			return i
		}
	}
	return -1
}

func ToInt64Array(array interface{}) []int64 {
	slice, len := ConvertToSlice(array)
	result := make([]int64, len)
	for i, item := range slice {
		switch item.(type) {
		case int64:
			result[i] = item.(int64)
		case uint64:
			result[i] = int64(item.(uint64))
		default:
			result[i] = cast.ToInt64(item)
		}
	}
	return result
}

func ToIntArray(array interface{}) []int {
	slice, len := ConvertToSlice(array)
	result := make([]int, len)
	for i, item := range slice {
		switch item.(type) {
		case int:
			result[i] = item.(int)
		default:
			result[i] = cast.ToInt(item)
		}
	}
	return result
}

func ToUInt64Array(array interface{}) []uint64 {
	slice, len := ConvertToSlice(array)
	result := make([]uint64, len)
	for i, item := range slice {
		result[i] = item.(uint64)
	}
	return result
}

func ToInterfaceArray(array interface{}) []interface{} {
	slice, len := ConvertToSlice(array)
	result := make([]interface{}, len)
	for i, item := range slice {
		result[i] = item
	}
	return result
}

// 将 interface 切片转为 string 切片
//
// interface 切片的元素类型支持 ObjectId 和 string
func InterfaceArrayToStringArray(array []interface{}) []string {
	result := make([]string, len(array))
	for i, item := range array {
		switch item.(type) {
		case bson.ObjectId:
			result[i] = (item.(bson.ObjectId)).Hex()
		case string:
			result[i] = item.(string)
		default:
			panic(fmt.Sprintf("not support type [%v] convert to string", reflect.TypeOf(item).Kind()))
		}
	}
	return result
}

// 将 array 转为 []interface{}，并返回 array 的长度
// 将 slice 和 array 以外的类型变量传入会引起 panic
func ConvertToSlice(array interface{}) ([]interface{}, int) {
	v := reflect.ValueOf(array)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		result := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			result[i] = v.Index(i).Interface()
		}
		return result, v.Len()
	case reflect.Ptr:
		return ConvertToSlice(v.Elem().Interface())
	default:
		panic(fmt.Sprintf("not support convert type[%s] to slice", v.Kind()))
	}
}

// 是否是个数组
func IsArray(array interface{}) bool {
	v := reflect.ValueOf(array)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		return true
	case reflect.Ptr:
		return IsArray(v.Elem().Interface())
	default:
		return false
	}
}

// RemoveValueItemsFromStructArrayByEqualFieldValue
//
// 从 structArray 中提取 key 字段的值，如果和 valueItems 中的一致，则删除 structArray 中的元素
// 返回切除后的切片和切片长度
// structArray 是一个结构体切片。valueItems 是一个切片，类型和 structArray 中单个元素以 key 为 fieldName 的字段类型一致。
func RemoveValueItemsFromStructArrayByEqualFieldValue(valueItems interface{}, structArray interface{}, key string) ([]interface{}, int) {
	slice, length := ConvertToSlice(structArray)

	isSlice := IsArray(valueItems)
	var valueItemsArray []interface{}
	if isSlice {
		valueItemsArray, _ = ConvertToSlice(valueItems)
	}

	for i := 0; i < length; i++ {
		value := GetValueByFiledName(slice[i], key)
		if (isSlice && (IndexOfInterfaceArray(value, valueItemsArray)) != -1) || (value == valueItems) {
			slice = append(slice[:i], slice[i+1:]...)
			i--
			length--
		}
	}
	return slice, length
}

// GetValueItemsFromStructArrayByEqualFieldValue
//
// 从 structArray 中提取 key 字段的值，如果和 valueItems 中的一致，则提取 structArray 中的元素
// 返回切除后的切片和切片长度
// structArray 是一个结构体切片。valueItems 是一个切片，类型和 structArray 中单个元素以 key 为 fieldName 的字段类型一致。
func GetValueItemsFromStructArrayByEqualFieldValue(valueItems interface{}, structArray interface{}, key string) ([]interface{}, int) {
	slice, length := ConvertToSlice(structArray)
	var result []interface{}
	isSlice := IsArray(valueItems)
	var valueItemsArray []interface{}
	if isSlice {
		valueItemsArray, _ = ConvertToSlice(valueItems)
	}

	for i := 0; i < length; i++ {
		value := GetValueByFiledName(slice[i], key)
		if (isSlice && (IndexOfInterfaceArray(value, valueItemsArray)) != -1) || (value == valueItems) {
			result = append(result, slice[i])
		}
	}
	return result, len(result)
}

// 从 valueArray 中去除包含与 structItems 中 key 字段值一致的元素
//
//	structItems 结构体切片
//	valueArray  待处理的值类型切片
//	key	    结构体切片元素的字段名
//
// 返回切除后的 valueArray 和 valueArray 长度
//
// valueArray 中的元素类型需要和 structItems 中 fieldName 为 key 的字段类型一致
func RemoveStructItemsFromValueArrayByEqualFieldValue(structItems interface{}, valueArray interface{}, key string) ([]interface{}, int) {
	slice, length := ConvertToSlice(structItems)
	valueSlice, valueSliceLength := ConvertToSlice(valueArray)
	for i := 0; i < length; i++ {
		value := GetValueByFiledName(slice[i], key)
		if index := IndexOfInterfaceArray(value, valueSlice); index != -1 {
			valueSlice = append(valueSlice[:index], valueSlice[index+1:]...)
			valueSliceLength--
		}
	}
	return valueSlice, valueSliceLength
}

// 将多个数组组成一个
// todo @alomerry Wu 数组过多可能性能下降
func CombineArrays(arrays ...interface{}) []interface{} {
	var result []interface{}
	for i, array := range arrays {
		if IsEmpty(reflect.ValueOf(array)) {
			continue
		}
		slice, _ := ConvertToSlice(array)
		if i == 0 {
			result = slice
		} else {
			result = append(result, slice...)
		}
	}
	return result
}

// 将数组按某个子数组字段展开
func Unwind(key string, array interface{}) interface{} {
	var result []interface{}
	slice := reflect.ValueOf(array)
	if !IsArray(array) {
		panic(fmt.Sprintf("not support convert type[%s] to slice", slice.Kind()))
	}
	for i := 0; i < slice.Len(); i++ {
		item := slice.Index(i)
		// 获取切片字段
		arrayField := GetValueByFiledName(item.Interface(), key)
		if arrayField == nil {
			result = append(result, item.Interface())
			continue
		}
		// 分解切片字段成多个结构体
		undividedArray, length := ConvertToSlice(arrayField)
		if len(undividedArray) == 0 {
			result = append(result, item.Interface())
			continue
		}
		dividedArray := make([]interface{}, length)
		for i := range undividedArray {
			tmp := item
			setValue := reflect.MakeSlice(reflect.TypeOf(arrayField), 1, 1)
			setValue.Index(0).Set(reflect.ValueOf(undividedArray[i]))
			tmp.FieldByName(key).Set(setValue)
			dividedArray[i] = tmp.Interface()
		}
		result = append(result, dividedArray...)
	}
	return result
}

// 根据给定某字段的值的数组排序
func SortArrayByFieldValues(field string, fieldValues interface{}, arr interface{}) {
	fieldOrderMap := map[interface{}]int{}
	tempSlice, _ := ConvertToSlice(arr)
	tempFieldValues, _ := ConvertToSlice(fieldValues)
	for i, v := range tempFieldValues {
		fieldOrderMap[v] = i
	}
	sort.SliceStable(arr, func(j, i int) bool {
		// 这里必须要重新给 tempSlice 赋值刷新顺序
		tempSlice, _ = ConvertToSlice(arr)
		return fieldOrderMap[GetValueByFiledName(tempSlice[j], field)] < fieldOrderMap[GetValueByFiledName(tempSlice[i], field)]
	})
}

// 根据给定的 uint64 数组中的值，按比例划分 value 的值，返回 values 对应下标划分到的值列表
//
// 如果按比例划分后有余数，剩余的数会按照 values 中的值由大到小的顺序，依次分配余数
func ScaleUint64ByValues(value uint64, values []uint64) []uint64 {
	var valuesSum uint64
	// 用于保存倒序排列后的 values 项和下标
	var sortedValuesWithIndex [][]uint64
	for i, v := range values {
		valuesSum += v
		sortedValuesWithIndex = append(sortedValuesWithIndex, []uint64{uint64(i), v})
	}
	// 根据 values 的项由大到小进行排序
	sort.SliceStable(sortedValuesWithIndex, func(i, j int) bool {
		return sortedValuesWithIndex[i][1] > sortedValuesWithIndex[j][1]
	})

	result := make([]uint64, len(values))
	if valuesSum == 0 {
		return result
	}
	remainingValue := value
	for i, v := range values {
		result[i] = value * v / valuesSum
		remainingValue -= result[i]
	}
	if remainingValue == 0 {
		return result
	}
	// 按比例划分一轮以后，remainingValue 一定不足以给每个 result 项再加 1
	// 只需要根据 values 中的值由大到小的顺序，依次给 result 对应下标的值加 1 直到 remainingValue 为 0 即可
	for _, valueWithIndex := range sortedValuesWithIndex {
		result[valueWithIndex[0]] += 1
		remainingValue -= 1
		if remainingValue == 0 {
			break
		}
	}

	return result
}

// TODO enhance
func CallAndCombineArray(funcName string, array []any) []interface{} {
	fieldValues := make([]interface{}, 0, len(array))
	for i := range array {
		fieldValues = append(fieldValues, reflect.ValueOf(array[i]).MethodByName(funcName).Call(make([]reflect.Value, 0))[0].Interface())
	}
	return fieldValues
}
