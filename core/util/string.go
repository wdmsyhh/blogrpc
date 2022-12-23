package util

import (
	"blogrpc/core/extension/bson"
	"blogrpc/core/util/algorithm"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// todo optimize by exkmp
func LastIndexOf(str, substr string) int {
	n := len(substr)
	switch {
	case n == 0:
		return 0
	case n == len(str):
		if substr == str {
			return 0
		}
		return -1
	case n > len(str):
		return -1
	}
	if strings.Index(str, substr) == -1 {
		return -1
	}
	lastIndex := 0
	for {
		index := strings.Index(str, substr)
		if index != -1 {
			lastIndex += index
			str = str[index+len(substr):]
			lastIndex += len(substr)
		} else {
			break
		}
	}
	return lastIndex - len(substr)
}

// 判断 array 数组是否包含 element
func ContainsString(array *[]string, element string) bool {
	return IndexOfArray(element, array) >= 0
}

func ContainsAll(array *[]string, element *[]string) bool {
	for _, item := range *element {
		if IndexOfArray(item, array) < 0 {
			return false
		}
	}
	return true
}

func FindStringSubmatch(matcher *regexp.Regexp, s string) string {
	if matcher == nil {
		return ""
	}
	temp := matcher.FindStringSubmatch(s)
	if len(temp) >= 2 {
		return temp[1]
	}
	return ""
}

func Remove(array *[]string, element string) {
	for i := 0; i < len(*array); i++ {
		if (*array)[i] == element {
			*array = append((*array)[:i], (*array)[i+1:]...)
			i--
		}
	}
}

func RemoveObjectId(array *[]bson.ObjectId, element bson.ObjectId) {
	for i := 0; i < len(*array); i++ {
		if (*array)[i].Hex() == element.Hex() {
			*array = append((*array)[:i], (*array)[i+1:]...)
			i--
		}
	}
}

func Removes(array *[]string, elements []string) {
	for _, str := range elements {
		Remove(array, str)
	}
}

func UppercaseFirst(word string) string {
	length := len(word)
	if length == 0 {
		return ""
	}
	remaining := word[1:]
	first := strings.ToUpper(string(word[0]))
	return strings.Join([]string{first, remaining}, "")
}

// Float64MultiHundred 仅用于计算部分旧字段使用 float 存储金额时转换成分时的精度问题
func Float64MultiHundred(amount float64) int64 {
	var negative int64 = 1
	if amount < 0 {
		negative = int64(-1)
	}
	amount = math.Abs(amount)
	amountStr := strings.Split(strconv.FormatFloat(amount, 'f', 2, 64), ".")
	result, _ := strconv.ParseInt(amountStr[0], 10, 64)
	result *= 100
	if len(amountStr) == 1 {
		return result * negative
	}
	for i, bit := 0, 10; i < len(amountStr[1]); i++ {
		result = result + int64(int(amountStr[1][i]-'0')*bit)
		bit /= 10
	}
	return result * negative
}

func StringArrayEqual(a, b *[]string) bool {
	if len(*a) != len(*b) {
		return false
	}
	ss := algorithm.CSet.InstanceFromStringSlice(a)
	for i := range *b {
		if !ss.Has((*b)[i]) {
			return false
		}
	}
	return true
}

func ToString(value interface{}) string {
	switch value.(type) {
	case string:
		return value.(string)
	case bson.ObjectId:
		return value.(bson.ObjectId).Hex()
	case primitive.ObjectID:
		return value.(primitive.ObjectID).Hex()
	default:
		return cast.ToString(value.(uint64))
	}
}

func IsAllNumber(str string) bool {
	for i := range str {
		if str[i] > '9' && str[i] < '0' {
			return false
		}
	}
	return true
}

func ReplaceEmoji(s string, r string) string {
	ret := ""
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if len(string(rs[i])) == 4 {
			ret = fmt.Sprintf("%s%s", ret, r)
		} else {
			ret = fmt.Sprintf("%s%s", ret, string(rs[i]))
		}
	}
	return ret
}

func NewUUIDWithServiceName() string {
	return fmt.Sprintf("%s_%s", os.Getenv("K8S_SERVICE_NAME"), uuid.NewV4().String())
}
