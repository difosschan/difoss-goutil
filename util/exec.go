package util

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const errFinished = "os: process already finished"

type Command struct {
	Cmd  string
	Args []string
}

type ExecResult struct {
	Stdout     string
	Stderr     string
	ExitStatus int
	Error      error
}

func NewCommand(cmd string, arg ...string) *Command {
	return &Command{
		Cmd:  cmd,
		Args: arg,
	}
}

func ExecCommand(name string, arg ...string) *ExecResult {
	cmd := exec.Command(name, arg...)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	e := cmd.Run()
	stdout := strings.TrimSpace(outBuf.String())
	stderr := strings.TrimSpace(errBuf.String())
	return NewExecResult(stdout, stderr, e)
}

func NewExecResult(stdout, stderr string, execErr error) *ExecResult {
	result := &ExecResult{
		Stdout: stdout,
		Stderr: stderr,
		Error:  execErr,
	}
	if execErr != nil {
		reg := regexp.MustCompile(`exit status (\d+)`)
		m := reg.FindStringSubmatch(execErr.Error())
		if len(m) > 0 {
			exitRet, err := strconv.ParseInt(m[1], 10, 64)
			if err != nil {
				result.ExitStatus = -1
				return result
			}
			result.ExitStatus = int(exitRet)
			result.Error = errors.New(stderr)
			return result
		}
		result.ExitStatus = -2
		return result
	}
	return result
}

func ExecCommandContext(ctx context.Context,
	name string, arg ...string) *ExecResult {

	cmd := exec.CommandContext(ctx, name, arg...)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	go func() {
		select {
		case <-ctx.Done():
			if err := cmd.Process.Signal(os.Interrupt); err != nil {
				if err.Error() == errFinished {
					return
				}
				fmt.Println("error when kill:", err.Error())
			}
		}
	}()

	err := cmd.Start()
	if err == nil {
		err = cmd.Wait()
	}

	return NewExecResult(strings.TrimSpace(outBuf.String()), strings.TrimSpace(errBuf.String()), err)
}

func FindByRegexp(in, pattern string) (match []string) {
	lines := strings.Split(in, "\n")
	var reg = regexp.MustCompile(pattern)
	for _, line := range lines {
		if m := reg.FindStringSubmatch(line); len(m) > 0 {
			return m
		}
	}
	return
}

func ExecCommandAndFind(cmd Command, pattern string) (exitStatus int, results []string, err error) {
	r := ExecCommand(cmd.Cmd, cmd.Args...)
	if r.Error != nil {
		return r.ExitStatus, results, r.Error
	}
	if len(pattern) != 0 {
		match := FindByRegexp(r.Stdout, pattern)
		results = match[1:]
	}
	return 0, results, nil
}
