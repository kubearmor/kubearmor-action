// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package os

import (
	"os"
)

func IsFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

func RemoceFile(fileName string) error {
	if IsFileExist(fileName) {
		return os.Remove(fileName)
	}
	return nil
}
