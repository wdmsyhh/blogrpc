package main

import (
	"fmt"
	"strings"
)

func main() {
	str := "a b c     d   "
	fs := strings.Fields(str)
	fmt.Println(fs)
}