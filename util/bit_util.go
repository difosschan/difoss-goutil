package util

import (
	"fmt"
	"strings"
)

func IsByteVisible(c byte) bool {
	return c > 31 && c < 128
}

func DumpHex(bs []byte, wordsPerLine ...int) (s string) {
	bsLen := len(bs)
	pl := 16
	if len(wordsPerLine) > 0 {
		pl = wordsPerLine[0]
	}

	if bsLen == 0 || pl <= 0 {
		return
	}
	s = fmt.Sprintf("(len = %d)\n", bsLen)
	lineCnt := (bsLen + pl - 1) / pl
	for i := 0; i < lineCnt; i++ {
		var (
			hexVec = make([]string, pl)
			visVec = make([]string, pl)
		)
		for j := 0; j < pl; j++ {
			var (
				ch    byte
				isIn  bool
				chIdx = i*pl + j
			)
			if chIdx < bsLen {
				ch = bs[chIdx]
				isIn = true
			}
			if isIn {
				hexVec[j] = fmt.Sprintf("%02X", ch&0xFF)
				if !IsByteVisible(ch) {
					ch = '.'
				}
				visVec[j] = fmt.Sprintf("%c", ch)
			} else {
				hexVec[j] = "  "
				visVec[j] = " "
			}
		}
		s += strings.Join(hexVec, " ") + "  " + strings.Join(visVec, "") + "\n"
	}
	return
}
