// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package utils

import (
	"fmt"
	"net/url"
	"os"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/kubearmor/kubearmor-action/utils/urlfile"
)

// Retry retries the given action for the given number of times
func Retry(tryTimes int, trySleepTime time.Duration, action func() error) error {
	var err error
	for i := 0; i < tryTimes; i++ {
		err = action()
		if err == nil {
			return nil
		}

		time.Sleep(trySleepTime * time.Duration(2*i+1))
	}
	return fmt.Errorf("retry action timeout: %v", err)
}

// GetUUID returns a UUID string
func GetUUID() string {
	// generate a new UUID
	id := uuid.New()

	// print the UUID
	return id.String()
}

// RemoveDuplication removes duplication from a string slice
func RemoveDuplication(arr []string) []string {
	length := len(arr)
	if length == 0 {
		// return empty slice
		return arr
	}
	// sort the slice
	sort.Strings(arr)
	// j is the index of the last unique element
	j := 0
	for i := 1; i < length; i++ {
		// if the current element is different from the last unique element
		if arr[i] != arr[j] {
			// increment j
			j++
			// if j is less than i
			if j < i {
				swap(arr, i, j)
			}
		}
	}
	// return the slice up to the last unique element
	return arr[:j+1]
}

// swap swaps two elements in a string slice
func swap(arr []string, a, b int) {
	arr[a], arr[b] = arr[b], arr[a]
}

// ReadFile reads the file from the given address and returns it as a []byte array.
// It can handle both remote urls and local paths.
func ReadFile(address string) ([]byte, error) {
	// Check if the address is a url
	if CheckIsURL(address) {
		// If the scheme is http or https, use http.Get to read the file
		data, err := urlfile.ReadJSONFromURL(address)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	// If the scheme is not http or https, assume it is a local path and use os.Open to read the file
	// os.ReadFile is a function that takes a file name and returns its content as a byte array and an error value
	data, err := os.ReadFile(address) // #nosec
	if err != nil {
		return nil, err
	}
	return data, nil
}

// CheckIsURL checks if the given address is a url
func CheckIsURL(address string) bool {
	// Parse the address as a url
	u, err := url.Parse(address)
	if err != nil {
		return false
	}
	// Check the scheme of the url
	if u.Scheme == "http" || u.Scheme == "https" {
		return true
	}
	return false
}
