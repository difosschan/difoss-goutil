package util

import (
	"fmt"
	"strings"
)

var replaceStr string

func init() {
	replaceStr = "*"
}

func HideString(s string, beginPosShown int, endPosShownOrNone ...int) string {
	l := len(s)
	var preHidden, sufHidden string

	begin, end := PythonLikePos(l, beginPosShown, endPosShownOrNone...)

	if beginPosShown > 0 {
		preHidden = strings.Repeat(replaceStr, beginPosShown)
	}
	if end < l {
		sufHidden = strings.Repeat(replaceStr, l-end)
	}

	return fmt.Sprintf("%s%s%s", preHidden, s[begin:end], sufHidden)
}

func ResetHiddenReplacer(r string) {
	replaceStr = r
}

func PythonLikeSlice(s string, beg int, endOrNone ...int) string {
	begin, end := PythonLikePos(len(s), beg, endOrNone...)
	return s[begin:end]
}

func PythonLikePos(length, beg int, endOrNone ...int) (begin, end int) {
	end = length
	if len(endOrNone) > 0 {
		end = transferPythonPosition(length, endOrNone[0])
	}
	begin = transferPythonPosition(length, beg)
	if begin > end {
		return begin, begin
	}
	return
}

func transferPythonPosition(length, posNav int) (pos int) {
	if posNav > length {
		return length
	}
	if posNav >= 0 {
		return posNav
	}
	pos = length + posNav
	if pos > 0 {
		return
	}
	return 0
}
