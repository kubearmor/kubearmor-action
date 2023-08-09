// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package os

import (
	"os"
	"path/filepath"
)

// FileWriter is an interface for writing files
type FileWriter interface {
	WriteFile(content []byte) error
}

// fileWriter is an implementation of FileWriter
type fileWriter struct {
	fileName string
}

// WriteFile writes contents to a file
func (c fileWriter) WriteFile(content []byte) error {
	dir := filepath.Dir(c.fileName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0750); err != nil {
			return err
		}
	}
	return os.WriteFile(c.fileName, content, 0644) // #nosec
}

// NewFileWriter returns a new FileWriter
func NewFileWriter(fileName string) FileWriter {
	return fileWriter{
		fileName: fileName,
	}
}
