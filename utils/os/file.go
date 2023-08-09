// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package os

import (
	"os"
)

// IsFileExist returns true if a file exists
func IsFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

// RemoveFile removes a file
func RemoveFile(fileName string) error {
	if IsFileExist(fileName) {
		return os.Remove(fileName)
	}
	return nil
}
