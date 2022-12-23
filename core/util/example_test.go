package util

import (
	"blogrpc/core/extension/bson"
	"fmt"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func ExampleConvertStructExistBasicFieldsToMap() {
	element := struct {
		Id      int
		Age     int
		Name    *string
		Married bool
	}{11, 18, nil, false}

	opt := ConvertStructExistBasicFieldsToMapOption{
		NeedLowercaseFirst: false,
		IsBoolFalseAsEmpty: false,
		HasId:              false,
		IgnoreCreatedAt:    true,
	}
	fmt.Println(ConvertStructExistBasicFieldsToMap(element, opt))
	opt.NeedLowercaseFirst = true
	fmt.Println(ConvertStructExistBasicFieldsToMap(element, opt))
	opt.NeedLowercaseFirst = true
	opt.IsBoolFalseAsEmpty = true
	opt.HasId = false
	fmt.Println(ConvertStructExistBasicFieldsToMap(element, opt))
	opt.NeedLowercaseFirst = true
	opt.IsBoolFalseAsEmpty = true
	opt.HasId = true
	fmt.Println(ConvertStructExistBasicFieldsToMap(element, opt))
	// output:
	// map[Age:18 Id:11 Married:false]
	// map[age:18 id:11 married:false]
	// map[age:18 id:11]
	// map[_id:11 age:18]
}

func ExampleContainsString() {
	array := []string{"abc", "edg", "?"}

	fmt.Println(ContainsString(&array, "?!"))
	fmt.Println(ContainsString(&array, "edg"))
	// output:
	// false
	// true
}

func ExampleContainsInt64() {
	array := []int64{1, 2, 3, -1}

	fmt.Println(ContainsInt64(&array, -2))
	fmt.Println(ContainsInt64(&array, int64(uint64(3))))
	// output:
	// false
	// true
}

func ExampleConvertToSlice() {
	array := []string{"abc", "edg", "?"}
	fmt.Println(ConvertToSlice(array))
	// output:
	// [abc edg ?] 3
}

func ExampleIsArray() {
	array := []interface{}{"name", "Mike", "age", 18}

	fmt.Println(IsArray(array))
	fmt.Println(IsArray(&array))
	// output:
	// true
	// true
}

func ExampleGetArraysIntersection() {
	type t1 struct {
		Id   int
		Name string
	}
	type t2 struct {
		Name string
	}
	array1, array2, array3, array4 :=
		[]t1{
			{1, "Red"},
			{2, "Green"},
			{3, "Yellow"},
			{4, "Orange"},
		},
		[]t1{
			{1, "Red"},
			{2, "Green"},
			{31, "Yellow"},
			{4, "Orange"},
		},
		[]t1{
			{1, "Red"},
			{21, "Green"},
			{31, "Yellow"},
			{4, "Orange"},
		},
		[]t2{
			{"Red"},
			{"Green"},
			{"Yellow"},
			{"Orange"},
		}

	res := GetArraysIntersection(array1, array2, array3)
	sort.Slice(res, func(i, j int) bool {
		return res[i].(t1).Id < res[j].(t1).Id
	})
	fmt.Println(res)
	fmt.Println(GetArraysIntersection(array2, array4))
	// output:
	// [{1 Red} {4 Orange}]
	// []
}

func TestGetArraysIntersection(t *testing.T) {
	type t1 struct {
		Id   int
		Name string
	}
	type t2 struct {
		Name string
	}
	array1, array2, array3 :=
		[]t1{
			{1, "Red"},
			{2, "Green"},
			{1, "Red"},
			{3, "Yellow"},
			{4, "Orange"},
		},
		[]t1{
			{1, "Red"},
			{2, "Green"},
			{2, "Green"},
			{31, "Yellow"},
			{4, "Orange"},
		},
		[]t1{
			{1, "Red"},
			{21, "Green"},
			{31, "Yellow"},
			{4, "Orange"},
		}
	res := GetArraysIntersection(array1, array2, array3)
	sort.Slice(res, func(i, j int) bool {
		return res[i].(t1).Id < res[j].(t1).Id
	})
	assert.True(t, reflect.DeepEqual(res[0], t1{1, "Red"}))
	assert.True(t, reflect.DeepEqual(res[1], t1{4, "Orange"}))
	assert.Equal(t, 2, len(res))
}

func BenchmarkGetArraysIntersection(b *testing.B) {
	type t1 struct {
		Id   int
		Name string
	}
	result := []interface{}{}
	for i := 0; i < 100; i++ {
		t := []t1{}
		for j := 0; j < 100; j++ {
			t = append(t, t1{
				Id:   rand.Intn(10),
				Name: cast.ToString(rand.Intn(50)),
			})
		}
		result = append(result, t)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b.ReportAllocs()
		GetArraysIntersection(result...)
	}
}

func ExampleExtractArrayField() {
	array := []struct {
		Id   int
		Name string
	}{
		{1, "Red"},
		{2, "Green"},
		{3, "Yellow"},
		{4, "Orange"},
	}

	mapper := []map[string]string{
		{"Id": "Red", "Age": "18"},
		{"Id": "Green", "Age": "21"},
		{"Id": "Yellow", "Age": "19"},
		{"Id": "Orange", "Age": "37"},
	}

	fmt.Println(ExtractArrayField("Id", array))
	fmt.Println(ExtractArrayField("Id", mapper))
	// output:
	// [1 2 3 4]
	// [Red Green Yellow Orange]
}

func ExampleGetFieldsInAndNotInElement() {
	fields := []string{"Name", "Age", "Location"}
	element1 := struct {
		Name   string
		Age    int
		Gender string
	}{"Mike", 18, "male"}
	element2 := map[string]interface{}{"Name": ""}
	element3 := []string{"Name", "Location"}

	fmt.Println(GetFieldsInAndNotInElement(fields, element1))
	fmt.Println(GetFieldsInAndNotInElement(fields, element2))
	fmt.Println(GetFieldsInAndNotInElement(fields, element3))
	// output:
	// [Name Age] [Location]
	// [Name] [Age Location]
	// [Name Location] [Age]
}

func ExampleGetValueByFiledName() {
	object := struct {
		Gender string
		Age    int
		Name   struct {
			FirstName  string
			SecondName string
		}
	}{
		"female",
		0,
		struct {
			FirstName  string
			SecondName string
		}{"Omega", "Mike"},
	}
	mapper := map[string]interface{}{"city": "shanghai"}

	fmt.Println(GetValueByFiledName(object, "Gender"))
	fmt.Println(GetValueByFiledName(object, "Name.FirstName"))
	fmt.Println(GetValueByFiledName(mapper, "city"))
	// output:
	// female
	// Omega
	// shanghai
}

func ExampleIndexOfArray() {
	array1 := []interface{}{bson.NewObjectId(), "hello"}

	fmt.Println(IndexOfArray("hello", array1))
	// output:
	// 1
}

func ExampleInterfaceArrayToStringArray() {
	array1 := []interface{}{bson.ObjectIdHex("5f48724cfcfbd7df4e8c163d"), "hello"}

	fmt.Println(InterfaceArrayToStringArray(array1))
	// output:
	// [5f48724cfcfbd7df4e8c163d hello]
}

func ExampleIsContains() {
	array := []string{"hello", "red"}
	object := struct {
		Name string
		Age  int
	}{"Mike", 18}
	mapper := map[string]interface{}{"city": "shanghai"}

	fmt.Println(IsContains("red", array))
	fmt.Println(IsContains("Name", object))
	fmt.Println(IsContains("Gender", object))
	fmt.Println(IsContains("city", mapper))
	fmt.Println(IsContains("country", mapper))
	// output:
	// true
	// true
	// false
	// true
	// false
}

func ExampleIsInMap() {
	mapper := map[interface{}]interface{}{"name": "Mike", "age": 18}

	fmt.Println(IsInMap("Name", mapper))
	fmt.Println(IsInMap("name", mapper))
	// output:
	// false
	// true
}

func ExampleIsStructHasField() {
	element := struct {
		Name string
		Age  int
	}{"Mike", 18}

	fmt.Println(IsStructHasField("Name", element))
	fmt.Println(IsStructHasField("Age", element))
	fmt.Println(IsStructHasField("Intro", element))
	// output:
	// true
	// true
	// false
}

func ExampleLowercaseFirst() {
	fmt.Println(LowercaseFirst("Id"))
	// output:
	// id
}

func ExampleMakeMapper() {
	array := []struct {
		Name string
		Age  int
	}{{"Apple", 41}, {"Pencil", 28}, {"Green", 16}}

	fmt.Println(MakeMapper("Name", array))
	// output:
	// map[Apple:{Apple 41} Green:{Green 16} Pencil:{Pencil 28}]
}

func ExampleRemoveStructItemsFromValueArrayByEqualFieldValue() {
	structItems := []struct {
		Name string
		Age  int
	}{{"Apple", 41}, {"Pencil", 28}, {"Green", 16}}
	valueArray1 := []string{"Apple", "Mongo", "Peach"}
	valueArray2 := []int{22, 16, 48, 41}

	fmt.Println(RemoveStructItemsFromValueArrayByEqualFieldValue(structItems, valueArray1, "Name"))
	fmt.Println(RemoveStructItemsFromValueArrayByEqualFieldValue(structItems, valueArray2, "Age"))
	// output:
	// [Mongo Peach] 2
	// [22 48] 2

}

func ExampleRemoveValueItemsFromStructArrayByEqualFieldValue() {
	structArray := []struct {
		Name string
		Age  int
	}{{"Apple", 41}, {"Pencil", 28}, {"Green", 16}}
	valueItems1 := []string{"Apple", "Mongo", "Peach"}
	valueItems2 := []int{22, 16, 48, 41}

	fmt.Println(RemoveValueItemsFromStructArrayByEqualFieldValue(valueItems1, structArray, "Name"))
	fmt.Println(RemoveValueItemsFromStructArrayByEqualFieldValue(valueItems2, structArray, "Age"))

	// output:
	// [{Pencil 28} {Green 16}] 2
	// [{Pencil 28}] 1
}

func ExampleRemove() {
	array1 := []string{"aef", "bmk", "cdk", "dny"}
	s := "dny"

	Remove(&array1, s)
	fmt.Println(array1)
	// output:
	// [aef bmk cdk]
}

func ExampleRemoves() {
	array1 := []string{"aef", "bmk", "cdk", "dny"}
	array2 := []string{"aef", "cdk"}

	Removes(&array1, array2)
	fmt.Println(array1)
	// output:
	// [bmk dny]
}

func ExampleToInterfaceArray() {
	array1 := []string{"a", "b"}
	array2 := []bool{false, true}

	fmt.Println(ToInterfaceArray(array1))
	fmt.Println(ToInterfaceArray(array2))
	// output:
	// [a b]
	// [false true]
}

func ExampleToObjectIdArray() {
	array1 := []interface{}{"5f486642fcfbd7db0884d209", "5f486642fcfbd7db0884d20a"}
	array2 := []interface{}{bson.ObjectIdHex("5f4866f5fcfbd7db5499c22a"), bson.ObjectIdHex("5f4866f5fcfbd7db5499c22b")}

	fmt.Println(ToObjectIdArray(array1))
	fmt.Println(ToObjectIdArray(array2))
	// output:
	// [ObjectIdHex("5f486642fcfbd7db0884d209") ObjectIdHex("5f486642fcfbd7db0884d20a")]
	// [ObjectIdHex("5f4866f5fcfbd7db5499c22a") ObjectIdHex("5f4866f5fcfbd7db5499c22b")]

}

func ExampleToStringArray() {
	array1 := []interface{}{"a", "b", "c"}
	array2 := []interface{}{bson.ObjectIdHex("5f486642fcfbd7db0884d209"), bson.ObjectIdHex("5f486642fcfbd7db0884d20a")}

	fmt.Println(ToStringArray(array1))
	fmt.Println(ToStringArray(array2))
	// output:
	// [a b c]
	// [5f486642fcfbd7db0884d209 5f486642fcfbd7db0884d20a]
}

func ExampleUppercaseFirst() {
	fmt.Println(UppercaseFirst("hello"))
	// output:
	// Hello
}

func ExampleGetMonthMaxDay() {
	fmt.Println(GetMonthMaxDay(1996, 2))
	fmt.Println(GetMonthMaxDay(2004, 12))
	fmt.Println(GetMonthMaxDay(1997, 4))
	// output:
	// 29
	// 31
	// 30
}

func ExampleUnwind() {
	person := struct {
		Color []string
		Name  string
	}{
		Color: []string{"Red", "Yellow"},
		Name:  "Mike",
	}

	fmt.Println(Unwind("Color", []struct {
		Color []string
		Name  string
	}{person}))
	// output:
	// [{[Red] Mike} {[Yellow] Mike}]
}

func ExampleCombineArrays() {
	fmt.Println(CombineArrays([]string{"A", "B", "C"}, []int{1, 2, 4}))
	// output:
	// [A B C 1 2 4]
}

func ExampleGetValueItemsFromStructArrayByEqualFieldValue() {
	people := []struct {
		Name string
		Age  int
	}{
		{Name: "Mike", Age: 16},
		{Name: "Jack", Age: 18},
		{Name: "Misaki", Age: 26},
		{Name: "Peach", Age: 46},
	}

	nameKeys := []string{"Misaki", "Peach", "Invalid"}
	fmt.Println(GetValueItemsFromStructArrayByEqualFieldValue(nameKeys, people, "Name"))
	// output:
	// [{Misaki 26} {Peach 46}] 2
}

func ExampleGetLastMonthDateRange() {
	fmt.Println(GetLastMonthDateRange(2004, 12))
	fmt.Println(GetLastMonthDateRange(2020, 9))
	// output:
	// 2004-11-01 00:00:00 +0800 CST 2004-12-01 00:00:00 +0800 CST
	// 2020-08-01 00:00:00 +0800 CST 2020-09-01 00:00:00 +0800 CST
}

func ExampleMakeArrayMapper() {
	array := []struct {
		Name   string
		Gender string
	}{
		{Name: "Ame", Gender: "Male"},
		{Name: "Misaki", Gender: "Female"},
		{Name: "Peach", Gender: "Female"},
		{Name: "Mei", Gender: "Female"},
	}

	fmt.Println(MakeArrayMapper("Gender", array))
	// output:
	// map[Female:[{Misaki Female} {Peach Female} {Mei Female}] Male:[{Ame Male}]]
}

func ExampleMakeMapperByJudge() {
	array := []struct {
		Name   string
		Gender string
	}{
		{Name: "Ame", Gender: "Male"},
		{Name: "Misaki", Gender: "Female"},
		{Name: "Peach", Gender: "Female"},
		{Name: "Mei", Gender: "Female"},
	}

	fmt.Println(MakeMapperByJudge("Name", array, func(item struct {
		Name   string
		Gender string
	}) bool {
		switch item.Name {
		case "Ame", "Misaki", "Peach":
			return true
		default:
			return false
		}
	}))
	// output:
	// map[Ame:{Ame Male} Misaki:{Misaki Female} Peach:{Peach Female}]
}

func ExampleSeparateSlice() {
	array1 := []struct {
		Name   string
		Gender string
	}{
		{Name: "Ame", Gender: "Male"},
		{Name: "Ani", Gender: "Female"},
		{Name: "Misaki", Gender: "Female"},
		{Name: "Peach", Gender: "Female"},
		{Name: "Mei", Gender: "Female"},
	}

	array2 := []struct {
		Name   string
		Gender string
	}{}

	fmt.Println(SeparateSlice(array1, 1))
	fmt.Println(SeparateSlice(array1, 2))
	fmt.Println(SeparateSlice(array1, 3))
	fmt.Println(SeparateSlice(array1, 4))
	fmt.Println(SeparateSlice(array1, 10))
	fmt.Println(SeparateSlice(array2, 10))

	// output:
	// [[{Ame Male} {Ani Female} {Misaki Female} {Peach Female} {Mei Female}]]
	// [[{Ame Male} {Ani Female}] [{Misaki Female} {Peach Female} {Mei Female}]]
	// [[{Ame Male}] [{Ani Female}] [{Misaki Female} {Peach Female} {Mei Female}]]
	// [[{Ame Male}] [{Ani Female}] [{Misaki Female}] [{Peach Female} {Mei Female}]]
	// [[{Ame Male}] [{Ani Female}] [{Misaki Female}] [{Peach Female}] [{Mei Female}]]
	// []
}
