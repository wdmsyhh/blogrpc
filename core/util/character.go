package util

import (
	"bytes"
	"encoding/hex"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
)

// GBK2312编码表
type character struct {
	// 前闭后开
	zhStart string
	zhEnd   string
	enCap   string
	en      string
}

var zhAndEnTable = map[character]string{
	{
		zhStart: "b0a1",
		zhEnd:   "b0c5",
		en:      "a3e1",
		enCap:   "a3c1",
	}: "A",
	{
		zhStart: "b0c5",
		zhEnd:   "b2c1",
		en:      "a3e2",
		enCap:   "a3c2",
	}: "B",
	{
		zhStart: "b2c1",
		zhEnd:   "b4ee",
		en:      "a3e3",
		enCap:   "a3c3",
	}: "C",
	{
		zhStart: "b4ee",
		zhEnd:   "b6ea",
		en:      "a3e4",
		enCap:   "a3c4",
	}: "D",
	{
		zhStart: "b6ea",
		zhEnd:   "b7a2",
		en:      "a3e5",
		enCap:   "a3c5",
	}: "E",
	{
		zhStart: "b7a2",
		zhEnd:   "b8c1",
		en:      "a3e6",
		enCap:   "a3c6",
	}: "F",
	{
		zhStart: "b8c1",
		zhEnd:   "b9fe",
		en:      "a3e7",
		enCap:   "a3c7",
	}: "G",
	{
		zhStart: "b9fe",
		zhEnd:   "bbf7",
		en:      "a3e8",
		enCap:   "a3c8",
	}: "H",
	{
		zhStart: "-1",
		zhEnd:   "-1",
		en:      "a3e9",
		enCap:   "a3c9",
	}: "I",
	{
		zhStart: "bbf7",
		zhEnd:   "bfa6",
		en:      "a3ea",
		enCap:   "a3ca",
	}: "J",
	{
		zhStart: "bfa6",
		zhEnd:   "c0ac",
		en:      "a3eb",
		enCap:   "a3cb",
	}: "K",
	{
		zhStart: "c0ac",
		zhEnd:   "c2e8",
		en:      "a3ec",
		enCap:   "a3cc",
	}: "L",
	{
		zhStart: "c2e8",
		zhEnd:   "c4c3",
		en:      "a3ed",
		enCap:   "a3cd",
	}: "M",
	{
		zhStart: "c4c3",
		zhEnd:   "c5b6",
		en:      "a3ee",
		enCap:   "a3ce",
	}: "N",
	{
		zhStart: "c5b6",
		zhEnd:   "c5be",
		en:      "a3ef",
		enCap:   "a3cf",
	}: "O",
	{
		zhStart: "c5be",
		zhEnd:   "c6da",
		en:      "a3f0",
		enCap:   "a3d0",
	}: "P",
	{
		zhStart: "c6da",
		zhEnd:   "c8bb",
		en:      "a3f1",
		enCap:   "a3d1",
	}: "Q",
	{
		zhStart: "c8bb",
		zhEnd:   "c8f6",
		en:      "a3f2",
		enCap:   "a3d2",
	}: "R",
	{
		zhStart: "c8f6",
		zhEnd:   "cBfa",
		en:      "a3f3",
		enCap:   "a3d3",
	}: "S",
	{
		zhStart: "cbfa",
		zhEnd:   "cdda",
		en:      "a3f4",
		enCap:   "a3d4",
	}: "T",
	{
		zhStart: "-1",
		zhEnd:   "-1",
		en:      "a3f5",
		enCap:   "a3d5",
	}: "U",
	{
		zhStart: "-1",
		zhEnd:   "-1",
		en:      "a3f6",
		enCap:   "a3d6",
	}: "V",
	{
		zhStart: "cdda",
		zhEnd:   "cef4",
		en:      "a3f7",
		enCap:   "a3d7",
	}: "W",
	{
		zhStart: "cef4",
		zhEnd:   "d1b9",
		en:      "a3f8",
		enCap:   "a3d8",
	}: "X",
	{
		zhStart: "d1b9",
		zhEnd:   "d4d1",
		en:      "a3f9",
		enCap:   "a3d9",
	}: "Y",
	{
		zhStart: "d4d1",
		zhEnd:   "d7fa",
		en:      "a3fa",
		enCap:   "a3da",
	}: "Z",
}

func IsASCIILetter(ch rune) bool {
	if (ch <= 90 && ch >= 65) || (ch <= 122 && ch >= 97) {
		return true
	}
	return false
}

func GetFirstLetter(chinese string) string {

	firstCh, err := utf8ToGbk([]byte(string([]rune(chinese)[:1])))
	if err != nil {
		return ""
	}

	// 检查字母
	if IsASCIILetter([]rune(chinese)[:1][0]) {
		return UppercaseFirst(string([]rune(chinese)[:1]))
	}

	hexCode := hex.EncodeToString(firstCh)

	// 检查 GBK2312 编码表
	for k, v := range zhAndEnTable {
		if (hexCode >= k.zhStart && hexCode < k.zhEnd) || hexCode == k.en || hexCode == k.enCap {
			return v
		}
	}
	return ""
}

func utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
