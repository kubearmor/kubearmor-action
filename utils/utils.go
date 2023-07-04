// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package utils

import (
	"fmt"
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
