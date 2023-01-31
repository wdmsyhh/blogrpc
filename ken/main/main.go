package main

import (
	"fmt"
	"strings"
)

func main() {
	line := "Value string `protobuf:\"bytes,1,opt,name=value\" json:\"value,omitempty\"`"
	res := strings.TrimRight(strings.TrimRight(line, "\n"), " ")
	res = strings.Trim(res, "`")
	if strings.Contains(line, "json:") {
		substr := " json:\"" + "value" + ",omitempty\""
		res = strings.Replace(res, substr, "", -1)
	}
	fmt.Println(res)
}
