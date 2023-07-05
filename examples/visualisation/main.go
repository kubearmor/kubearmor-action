// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package main

import (
	"fmt"

	"github.com/kubearmor-action/common"
	visual "github.com/kubearmor-action/pkg/visualisation"
)

func main() {
	// jsonFile is the name of the JSON file that you shared with me
	jsonFile := common.GetWorkDir() + "/test/testdata/test-summary-data.json"
	err := visual.ConvertJsonToImage(jsonFile, "test.png")
	if err != nil {
		fmt.Println("Error:", err)
	}
}
