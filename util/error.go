package util

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type MyError struct {
	Error  error
	Reason string
}

func NewError(err error, reason ...string) *MyError {
	return &MyError{
		err,
		strings.Join(reason, " ,"),
	}
}

func Panic(err error, reason ...string) {
	panic(NewError(err, reason...))
}

func DealWithPanic(panicVal interface{}) *MyError {
	er := &MyError{}
	if e, ok := panicVal.(error); ok && e != nil {
		er.Reason = "ERROR"
		er.Error = e
	} else if errorReason, ok := panicVal.(*MyError); ok && errorReason != nil {
		er = errorReason
	} else {
		er.Reason = "PANIC"
		er.Error = errors.New(fmt.Sprintf(": %v", panicVal))
	}
	return er
}

func PrintStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	fmt.Printf("STACK: %s\n", string(buf[:n]))
}
