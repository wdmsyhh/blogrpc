package algorithm

import (
	"crypto/sha256"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"sort"
	"testing"
)

var (
	h = sha256.New()
)

func TestInstanceFromStructSlice1(t *testing.T) {
	type SimpleStruct struct {
		Name string
		Age  int
		Sex  bool
	}

	ss := []SimpleStruct{
		{"Misaki", 1, true},
		{"Mike", 2, false},
		{"Eleven", 3, true},
		{"Eleven", 3, true},
		{"Eleven", 4, true},
		{"Eleven", 4, false},
	}
	rss := CSet.InstanceFromStructSlice(ss).ToArray()
	sort.Slice(rss, func(i, j int) bool {
		return getSha256Hash(rss[i].(SimpleStruct)) > getSha256Hash(rss[j].(SimpleStruct))
	})

	assert.Equal(t, rss[0].(SimpleStruct), SimpleStruct{Name: "Eleven", Age: 4, Sex: true})
	assert.Equal(t, rss[1].(SimpleStruct), SimpleStruct{Name: "Mike", Age: 2, Sex: false})
	assert.Equal(t, rss[2].(SimpleStruct), SimpleStruct{Name: "Eleven", Age: 4, Sex: false})
	assert.Equal(t, rss[3].(SimpleStruct), SimpleStruct{Name: "Misaki", Age: 1, Sex: true})
	assert.Equal(t, rss[4].(SimpleStruct), SimpleStruct{Name: "Eleven", Age: 3, Sex: true})
}

func TestInstanceFromStructSlice2(t *testing.T) {
	type Detail struct {
		Phone string
		Hobby *[]string
	}
	type SimpleStruct struct {
		Name   string
		Detail Detail
	}
	ss := []SimpleStruct{
		{"E", Detail{Phone: "6"}},
		{"E", Detail{Phone: "67"}},
		{"E", Detail{Phone: "6", Hobby: &[]string{"R", "W"}}},
		{"E", Detail{Phone: "6", Hobby: &[]string{"R", "W", "RW"}}},
		{"E", Detail{Phone: "6", Hobby: &[]string{"R", "W"}}},
		{"E", Detail{Phone: "6", Hobby: &[]string{"W", "R"}}},
	}

	should := []SimpleStruct{
		{"E", Detail{Phone: "6"}},
		{"E", Detail{Phone: "67"}},
		{"E", Detail{Phone: "6", Hobby: &[]string{"R", "W", "RW"}}},
		{"E", Detail{Phone: "6", Hobby: &[]string{"W", "R"}}},
		{"E", Detail{Phone: "6", Hobby: &[]string{"R", "W"}}},
	}

	rss := CSet.InstanceFromStructSlice(ss).ToArray()
	sort.Slice(rss, func(i, j int) bool {
		return getSha256Hash(rss[i].(SimpleStruct)) > getSha256Hash(rss[j].(SimpleStruct))
	})

	for i := range rss {
		var exists bool
		for j := range should {
			if reflect.DeepEqual(rss[i], should[j]) {
				exists = true
				break
			}
		}
		assert.Truef(t, exists, fmt.Sprintf("%v, %#v", i, rss[i]))
	}
}

func TestInstanceFromStructSlice3(t *testing.T) {
	type Detail struct {
		Phone string
		Email string
	}

	type SimpleStruct struct {
		Name   string
		Age    int
		Detail Detail
		Hobby  *[]string
	}
	ss := SimpleStruct{"Eleven", 3, Detail{Phone: "123456", Email: "123456@163.com"}, nil}
	assert.True(t, hasPointerFiled(reflect.TypeOf(ss)))
	assert.True(t, !hasPointerFiled(reflect.TypeOf(ss.Detail)))
}

func TestInstanceFromStructSlice4(t *testing.T) {
	type Detail struct {
		Phone string
		Email string
	}

	type SimpleStruct struct {
		Name   string
		Age    int
		Detail Detail
		Hobby  *[]string
	}
	sss1 := SimpleStruct{"Eleven", 3, Detail{Phone: "123456", Email: "123456@163.com"}, nil}
	sss2 := SimpleStruct{"Eleven", 3, Detail{Phone: "123456", Email: "123456@163.com"}, nil}
	sss3 := SimpleStruct{"Eleven", 3, Detail{Phone: "123456", Email: "123456@163.com"}, &[]string{"Read"}}
	sss4 := SimpleStruct{"Eleven", 3, Detail{Phone: "123456", Email: "123456@163.com"}, &[]string{"Read"}}
	assert.True(t, sss1 == sss2)
	assert.True(t, sss1 != sss3)
	assert.True(t, sss3 != sss4)
}

func getSha256Hash(v interface{}) string {
	h.Reset()
	h.Write([]byte(fmt.Sprintf("%v", v)))
	return string(h.Sum(nil))
}
