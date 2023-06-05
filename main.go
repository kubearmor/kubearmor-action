// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package main

import (
	"fmt"

	"github.com/sethvargo/go-githubactions"
)

func main() {
	action := githubactions.New()
	oldAppName := action.GetInput("old-app-image-name")
	if oldAppName == "" {
		action.Fatalf("old-app-image-name cannot be empty")
	}
	newAppName := action.GetInput("new-app-image-name")
	if newAppName == "" {
		action.Fatalf("new-app-image-name cannot be empty")
	}
	filepath := action.GetInput("filepath")

	fmt.Printf("oldAppName: %v, newAppName: %v, filepath: %v", oldAppName, newAppName, filepath)
}
