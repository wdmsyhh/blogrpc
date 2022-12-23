package algorithm

import (
	"blogrpc/core/extension/bson"
	"fmt"
	"reflect"
)

type void struct{}

var CSet = Set{}

type Set struct {
	set map[interface{}]void
}

func (s *Set) Instance() *Set {
	ss := Set{make(map[interface{}]void)}
	return &ss
}

func (s *Set) InstanceFromStringSlice(slice *[]string) *Set {
	ss := CSet.Instance()
	for i := range *slice {
		ss = ss.Insert((*slice)[i])
	}
	return ss
}

// 注意：结构体中有 map、slice、channel 或 func 字段不支持此方法，指针数组中元素的顺序不同视为不一致
func (s *Set) InstanceFromStructSlice(array interface{}) *Set {
	slice, len := convertToSlice(array)
	ss := CSet.Instance()

	// 检查结构体是否包含指针
	if len > 0 && hasPointerFiled(reflect.TypeOf(slice[0])) {
		// 包含指针时使用 deepEqual 来插入，TODO 复杂度 n²，可优化
		var exists bool
		for i := range slice {
			exists = false
			for k := range ss.set {
				v := k
				if reflect.DeepEqual(v, slice[i]) {
					exists = true
					break
				}
			}
			if !exists {
				ss = ss.Insert(slice[i])
			}
		}
	} else {
		for i := range slice {
			ss = ss.Insert(slice[i])
		}
	}
	return ss
}

func (s *Set) InstanceFromObjectIdSlice(slice *[]bson.ObjectId) *Set {
	ss := CSet.Instance()
	for i := range *slice {
		ss = ss.Insert((*slice)[i])
	}
	return ss
}

func (s *Set) Clear() {
	for k := range s.set {
		delete(s.set, k)
	}
}

func (s *Set) Empty() bool {
	return len(s.set) == 0
}

func (s *Set) Size() uint {
	if s == nil {
		return 0
	}
	return uint(len(s.set))
}

func (s *Set) Insert(val interface{}) *Set {
	if s == nil {
		set := Set{make(map[interface{}]void)}
		s = &set
	}
	s.set[val] = void{}
	return s
}

func (s *Set) TryInsert(val interface{}) (*Set, bool) {
	// 是否不存在并插入成功
	var insertSuccess bool
	if s == nil {
		set := Set{make(map[interface{}]void)}
		s = &set // 无效的，务必提前创建 Set
		insertSuccess = true
	} else {
		_, ok := s.set[val]
		insertSuccess = !ok
		s.set[val] = void{}
	}

	return s, insertSuccess
}

func (s *Set) InsertAll(val ...interface{}) *Set {
	if s == nil {
		set := Set{make(map[interface{}]void)}
		s = &set
	}
	for i := range val {
		s.set[val[i]] = void{}
	}
	return s
}

func (s *Set) Remove(val interface{}) {
	delete(s.set, val)
}

func (s *Set) Has(val interface{}) bool {
	_, exists := s.set[val]
	return exists
}

func (s *Set) HasAnyItem(val ...interface{}) bool {
	for i := range val {
		if _, exists := s.set[val[i]]; exists {
			return true
		}
	}
	return false
}

func (s *Set) ToArray() []interface{} {
	result := make([]interface{}, s.Size())
	index := 0
	for k := range s.set {
		result[index] = k
		index++
	}
	return result
}

func (s *Set) Clone() *Set {
	result := Set{make(map[interface{}]void, s.Size())}
	reflect.Copy(reflect.ValueOf(s), reflect.ValueOf(&result))
	return &result
}

// 防止循环引用
func convertToSlice(array interface{}) ([]interface{}, int) {
	v := reflect.ValueOf(array)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		result := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			result[i] = v.Index(i).Interface()
		}
		return result, v.Len()
	case reflect.Ptr:
		return convertToSlice(v.Elem().Interface())
	default:
		panic(fmt.Sprintf("not support convert type[%s] to slice", v.Kind()))
	}
}

// 判断结构体是否包含指针字段，不判断 map、slice、channel 或 func
func hasPointerFiled(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Ptr:
		return true
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			if hasPointerFiled(t.Field(i).Type) {
				return true
			}
		}
	}
	return false
}
