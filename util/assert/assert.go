package assert

import (
	"fmt"
	"difoss-goutil/util"
	"testing"
)

var innerT *testing.T

func SetT(t *testing.T) {
	innerT = t
}

func OK(expr bool, fatalInfo string) {
	if ! expr {
		file, line, funcName, _ := util.GetFuncInfo(1)
		out := fmt.Sprintf("FAIL <%v:%d:%s> %s", file, line, funcName, fatalInfo)
		if innerT != nil {
			innerT.Fatal(out)
		} else {
			panic(fmt.Errorf(out))
		}
	}
}
