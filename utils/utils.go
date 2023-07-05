// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package utils

import (
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
)

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
