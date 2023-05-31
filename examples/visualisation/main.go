// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package main

import (
	"fmt"

	"github.com/kubearmor/kubearmor-action/common"
	visual "github.com/kubearmor/kubearmor-action/pkg/visualisation"
)

func main() {
	// jsonFile is the name of the JSON file that you shared with me
	old_jsonFile := common.GetWorkDir() + "/test/testdata/test.json"
	new_jsonFile := common.GetWorkDir() + "/test/testdata/test.json"

	// sd := visual.ParseSummaryData(jsonFile)
	// vnd := visual.ParseNetworkData(sd)
	// fmt.Println(vnd)
	appName := "cassandra"
	err := visual.ConvertNetworkJSONToImage(old_jsonFile, new_jsonFile, "test.png", "")
	if err != nil {
		fmt.Println("Network-Visualisation Error:", err)
	}
	err = visual.ConvertSysJSONToImage(new_jsonFile, "test2.png", appName)
	if err != nil {
		fmt.Println("System-Visualisation Error:", err)
	}
}
