// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package common

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

const (
	OldAppTemplateFilePath     = "./template/app-old.yaml"        // Path to the old app template yaml file
	NewAppTemplateFilePath     = "./template/app-new.yaml"        // Path to the new app template yaml file
	OldAppImagePlaceholderName = "old-app-image-name-placeholder" // The old app image placeholder name to be changed
	NewAppImagePlaceholderName = "new-app-image-name-placeholder" // The new app image placeholder name to be changed
	AppNamespace               = "app"                            // Namespace is the default namespace of the app
)

const (
	ROOT = "root"
)

var (
	StdOut = os.Stdout
	StdErr = os.Stderr
)

func DefaultKubeConfigDir() string {
	return filepath.Join(GetHomeDir(), ".kube")
}

func GetHomeDir() string {
	home, err := homedir.Dir()
	if err != nil {
		return "/root"
	}
	return home
}
