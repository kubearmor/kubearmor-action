// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package exec

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/zhy76/kubearmor-action/common"

	"github.com/sirupsen/logrus"
)

// Cmd runs a command
func Cmd(name string, args ...string) error {
	cmd := exec.Command(name, args[:]...) // #nosec
	cmd.Stdin = os.Stdin
	cmd.Stderr = common.StdErr
	cmd.Stdout = common.StdOut
	return cmd.Run()
}

// CmdOutput runs a command and returns its combined standard output and standard error.
func CmdOutput(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args[:]...) // #nosec
	return cmd.CombinedOutput()
}

// RunSimpleCmd runs a simple command
func RunSimpleCmd(cmd string) (string, error) {
	var result []byte
	result, err := exec.Command("/bin/bash", "-c", cmd).CombinedOutput() // #nosec
	if err != nil {
		logrus.Debugf("failed to execute command(%s): error(%v)", cmd, err)
	}
	return string(result), err
}

// CheckCmdIsExist checks if the command exists
func CheckCmdIsExist(cmd string) (string, bool) {
	cmd = fmt.Sprintf("type %s", cmd)
	out, err := RunSimpleCmd(cmd)
	if err != nil {
		return "", false
	}

	outSlice := strings.Split(out, "is")
	last := outSlice[len(outSlice)-1]

	if last != "" && !strings.Contains(last, "not found") {
		return strings.TrimSpace(last), true
	}
	return "", false
}

// GetCurrentUserName returns the current user name
func GetCurrentUserName() (string, error) {
	u, err := user.Current()
	return u.Username, err
}
