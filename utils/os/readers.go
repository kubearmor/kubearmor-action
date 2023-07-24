// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package os

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type FileReader interface {
	ReadLines() ([]string, error)
	ReadAll() ([]byte, error)
}

type fileReader struct {
	fileName string
}

func (r fileReader) ReadLines() ([]string, error) {
	var lines []string

	if _, err := os.Stat(r.fileName); err != nil || os.IsNotExist(err) {
		return nil, errors.New("no such file")
	}

	file, err := os.Open(filepath.Clean(r.fileName))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			logrus.Fatal("failed to close file")
		}
	}()
	br := bufio.NewReader(file)
	for {
		line, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		lines = append(lines, string(line))
	}
	return lines, nil
}

func (r fileReader) ReadAll() ([]byte, error) {
	if _, err := os.Stat(r.fileName); err != nil || os.IsNotExist(err) {
		return nil, errors.New("no such file")
	}

	file, err := os.Open(filepath.Clean(r.fileName))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			logrus.Errorf("failed to close file: %v", err)
		}
	}()

	content, err := os.ReadFile(filepath.Clean(r.fileName))
	if err != nil {
		return nil, err
	}

	return content, nil
}

func NewFileReader(fileName string) FileReader {
	return fileReader{
		fileName: fileName,
	}
}
