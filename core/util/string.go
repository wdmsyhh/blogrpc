package util

import "strings"

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
