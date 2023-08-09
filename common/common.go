// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package common

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

const (
	// LOCALHOST IP address
	LOCALHOST = "127.0.0.1"
)

var (
	// StdOut is a standard output
	StdOut = os.Stdout
	// StdErr is a standard error
	StdErr = os.Stderr
)

// DefaultKubeConfigDir returns the default kubeconfig directory
func DefaultKubeConfigDir() string {
	return filepath.Join(GetHomeDir(), ".kube")
}

// GetHomeDir returns the home directory
func GetHomeDir() string {
	home, err := homedir.Dir()
	if err != nil {
		return "/root"
	}
	return home
}

// GetWorkDir returns the working directorys
func GetWorkDir() string {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return pwd
}
